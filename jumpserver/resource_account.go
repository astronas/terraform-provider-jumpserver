package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const accountBasePath = "accounts/accounts/"

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		ReadContext:   resourceAccountRead,
		UpdateContext: resourceAccountUpdate,
		DeleteContext: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username for authentication. Enter null if no username is required. For AD accounts use username@domain.",
			},
			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "password",
				ValidateFunc: validation.StringInSlice([]string{
					"password", "ssh_key", "access_key", "token", "api_key",
				}, false),
				Description: "The type of secret: password, ssh_key, access_key, token, api_key.",
			},
			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The secret value (password, SSH key, etc.).",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase for SSH key (only used when secret_type is ssh_key).",
			},
			"asset": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the asset this account belongs to.",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Account template ID to create account from.",
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
				Description: "Whether the account is active.",
			},
			"su_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the account to switch from (su).",
			},
			"on_invalid": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "error",
				ValidateFunc: validation.StringInSlice([]string{
					"error", "skip", "update",
				}, false),
				Description: "Policy when account already exists: error, skip, update.",
			},
			"push_now": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to push the account to the asset immediately after creation.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the account.",
			},
		},
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	accountData := map[string]interface{}{
		"name":        d.Get("name").(string),
		"username":    d.Get("username").(string),
		"secret_type": d.Get("secret_type").(string),
		"asset":       d.Get("asset").(string),
		"privileged":  d.Get("privileged").(bool),
		"is_active":   d.Get("is_active").(bool),
		"on_invalid":  d.Get("on_invalid").(string),
		"push_now":    d.Get("push_now").(bool),
	}

	if v, ok := d.GetOk("secret"); ok {
		accountData["secret"] = v.(string)
	}
	if v, ok := d.GetOk("passphrase"); ok {
		accountData["passphrase"] = v.(string)
	}
	if v, ok := d.GetOk("su_from"); ok {
		accountData["su_from"] = v.(string)
	}
	if v, ok := d.GetOk("template"); ok {
		accountData["template"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		accountData["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", accountBasePath, accountData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating account: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving account ID after creation, response: %v", result)
	}

	return resourceAccountRead(ctx, d, m)
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", accountBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading account: %s", resp.Status)
	}

	var account map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, account, "name")
	setStringField(d, account, "username")
	setEnumField(d, account, "secret_type")
	setObjectIDField(d, account, "asset")
	setObjectIDField(d, account, "template")
	setBoolField(d, account, "privileged")
	setBoolField(d, account, "is_active")
	setObjectIDField(d, account, "su_from")
	setStringField(d, account, "comment")

	return diags
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	accountData := map[string]interface{}{
		"name":        d.Get("name").(string),
		"username":    d.Get("username").(string),
		"secret_type": d.Get("secret_type").(string),
		"asset":       d.Get("asset").(string),
		"privileged":  d.Get("privileged").(bool),
		"is_active":   d.Get("is_active").(bool),
	}

	if v, ok := d.GetOk("secret"); ok {
		accountData["secret"] = v.(string)
	}
	if v, ok := d.GetOk("passphrase"); ok {
		accountData["passphrase"] = v.(string)
	}
	if v, ok := d.GetOk("su_from"); ok {
		accountData["su_from"] = v.(string)
	}
	if v, ok := d.GetOk("template"); ok {
		accountData["template"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		accountData["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", accountBasePath+d.Id()+"/", accountData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating account: %s", resp.Status)
	}

	return resourceAccountRead(ctx, d, m)
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", accountBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting account: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
