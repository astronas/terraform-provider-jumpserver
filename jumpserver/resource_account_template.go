package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const accountTemplateBasePath = "accounts/account-templates/"

func resourceAccountTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountTemplateCreate,
		ReadContext:   resourceAccountTemplateRead,
		UpdateContext: resourceAccountTemplateUpdate,
		DeleteContext: resourceAccountTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account template.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username. For AD accounts use username@domain.",
			},
			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "password",
				ValidateFunc: validation.StringInSlice([]string{
					"password", "ssh_key",
				}, false),
				Description: "The type of secret: password or ssh_key.",
			},
			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The secret value (password or SSH key content).",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase for SSH key.",
			},
			"secret_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "specific",
				ValidateFunc: validation.StringInSlice([]string{
					"specific", "random",
				}, false),
				Description: "Secret strategy: specific (user-provided) or random (auto-generated).",
			},
			"password_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      16,
							ValidateFunc: validation.IntBetween(8, 36),
							Description:  "Password length (8-36).",
						},
						"lowercase": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include lowercase letters.",
						},
						"uppercase": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include uppercase letters.",
						},
						"digit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include digits.",
						},
						"symbol": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include symbols.",
						},
						"exclude_symbols": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Symbols to exclude from generated passwords.",
						},
					},
				},
				Description: "Password complexity rules (used when secret_strategy is random).",
			},
			"platforms": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of platform IDs to associate with this template.",
			},
			"su_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account template ID to switch from (su).",
			},
			"privileged": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is a privileged account.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the template is active.",
			},
			"auto_push": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically push account to assets.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the template.",
			},
		},
	}
}

func accountTemplateBuildPayload(d *schema.ResourceData) map[string]interface{} {
	data := map[string]interface{}{
		"name":            d.Get("name").(string),
		"secret_type":     d.Get("secret_type").(string),
		"secret_strategy": d.Get("secret_strategy").(string),
		"privileged":      d.Get("privileged").(bool),
		"is_active":       d.Get("is_active").(bool),
		"auto_push":       d.Get("auto_push").(bool),
		"comment":         d.Get("comment").(string),
	}

	if v, ok := d.GetOk("username"); ok {
		data["username"] = v.(string)
	}
	if v, ok := d.GetOk("secret"); ok {
		data["secret"] = v.(string)
	}
	if v, ok := d.GetOk("passphrase"); ok {
		data["passphrase"] = v.(string)
	}
	if v, ok := d.GetOk("su_from"); ok {
		data["su_from"] = v.(string)
	}
	if v, ok := d.GetOk("platforms"); ok {
		data["platforms"] = v.([]interface{})
	}

	if v, ok := d.GetOk("password_rules"); ok {
		rules := v.([]interface{})
		if len(rules) > 0 {
			r := rules[0].(map[string]interface{})
			data["password_rules"] = map[string]interface{}{
				"length":          r["length"],
				"lowercase":       r["lowercase"],
				"uppercase":       r["uppercase"],
				"digit":           r["digit"],
				"symbol":          r["symbol"],
				"exclude_symbols": r["exclude_symbols"],
			}
		}
	}

	return data
}

func resourceAccountTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := accountTemplateBuildPayload(d)

	resp, err := c.doRequest("POST", accountTemplateBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating account template: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving account template ID after creation, response: %v", result)
	}

	return resourceAccountTemplateRead(ctx, d, m)
}

func resourceAccountTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", accountTemplateBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading account template: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "username")
	setBoolField(d, result, "privileged")
	setBoolField(d, result, "is_active")
	setBoolField(d, result, "auto_push")
	setStringField(d, result, "comment")

	setEnumField(d, result, "secret_type")
	setEnumField(d, result, "secret_strategy")
	setObjectIDField(d, result, "su_from")
	setObjectIDsField(d, result, "platforms")

	if pr := readPasswordRules(result); pr != nil {
		d.Set("password_rules", pr)
	}

	return diags
}

func readPasswordRules(data map[string]interface{}) []interface{} {
	v, ok := data["password_rules"].(map[string]interface{})
	if !ok {
		return nil
	}
	rules := map[string]interface{}{
		"length":          16,
		"lowercase":       true,
		"uppercase":       true,
		"digit":           true,
		"symbol":          true,
		"exclude_symbols": "",
	}
	if l, ok := v["length"].(float64); ok {
		rules["length"] = int(l)
	}
	if l, ok := v["lowercase"].(bool); ok {
		rules["lowercase"] = l
	}
	if u, ok := v["uppercase"].(bool); ok {
		rules["uppercase"] = u
	}
	if d2, ok := v["digit"].(bool); ok {
		rules["digit"] = d2
	}
	if s, ok := v["symbol"].(bool); ok {
		rules["symbol"] = s
	}
	if e, ok := v["exclude_symbols"].(string); ok {
		rules["exclude_symbols"] = e
	}
	return []interface{}{rules}
}

func resourceAccountTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := accountTemplateBuildPayload(d)

	resp, err := c.doRequest("PUT", accountTemplateBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating account template: %s", resp.Status)
	}

	return resourceAccountTemplateRead(ctx, d, m)
}

func resourceAccountTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", accountTemplateBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting account template: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
