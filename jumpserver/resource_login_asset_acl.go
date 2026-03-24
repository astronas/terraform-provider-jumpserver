package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const loginAssetACLBasePath = "acls/login-asset-acls/"

func resourceLoginAssetACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoginAssetACLCreate,
		ReadContext:   resourceLoginAssetACLRead,
		UpdateContext: resourceLoginAssetACLUpdate,
		DeleteContext: resourceLoginAssetACLDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the login asset ACL.",
			},
			"users": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User filter as JSON (e.g. {\"type\":\"ids\",\"ids\":[\"uuid1\"]}).",
			},
			"assets": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Asset filter as JSON (e.g. {\"type\":\"ids\",\"ids\":[\"uuid1\"]}).",
			},
			"accounts": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of account usernames (e.g. [\"root\",\"@ALL\"]).",
			},
			"rules": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rules as JSON object (flexible structure for IP, time conditions).",
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validation.IntBetween(1, 100),
				Description:  "Priority (1-100). Lower values match first.",
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "reject",
				ValidateFunc: validation.StringInSlice([]string{
					"reject", "accept", "review", "notice",
				}, false),
				Description: "Action when matched: reject, accept, review, notice.",
			},
			"reviewers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of reviewer user IDs (required when action is review).",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the ACL is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the ACL.",
			},
		},
	}
}

func resourceLoginAssetACLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}
	var assets interface{}
	if err := json.Unmarshal([]byte(d.Get("assets").(string)), &assets); err != nil {
		return diag.Errorf("assets must be valid JSON: %s", err)
	}
	var rules interface{}
	if err := json.Unmarshal([]byte(d.Get("rules").(string)), &rules); err != nil {
		return diag.Errorf("rules must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"users":     users,
		"assets":    assets,
		"accounts":  d.Get("accounts").([]interface{}),
		"rules":     rules,
		"priority":  d.Get("priority").(int),
		"action":    d.Get("action").(string),
		"is_active": d.Get("is_active").(bool),
	}

	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", loginAssetACLBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating login asset ACL: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving login asset ACL ID after creation, response: %v", result)
	}

	return resourceLoginAssetACLRead(ctx, d, m)
}

func resourceLoginAssetACLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", loginAssetACLBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading login asset ACL: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setIntField(d, result, "priority")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")

	if v, ok := result["users"]; ok {
		b, _ := json.Marshal(v)
		d.Set("users", string(b))
	}
	if v, ok := result["assets"]; ok {
		b, _ := json.Marshal(v)
		d.Set("assets", string(b))
	}
	if v, ok := result["accounts"].([]interface{}); ok {
		d.Set("accounts", v)
	}
	if v, ok := result["rules"]; ok {
		b, _ := json.Marshal(v)
		d.Set("rules", string(b))
	}

	setEnumField(d, result, "action")
	setObjectIDsField(d, result, "reviewers")

	return diags
}

func resourceLoginAssetACLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}
	var assets interface{}
	if err := json.Unmarshal([]byte(d.Get("assets").(string)), &assets); err != nil {
		return diag.Errorf("assets must be valid JSON: %s", err)
	}
	var rules interface{}
	if err := json.Unmarshal([]byte(d.Get("rules").(string)), &rules); err != nil {
		return diag.Errorf("rules must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"users":     users,
		"assets":    assets,
		"accounts":  d.Get("accounts").([]interface{}),
		"rules":     rules,
		"priority":  d.Get("priority").(int),
		"action":    d.Get("action").(string),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}

	resp, err := c.doRequest("PUT", loginAssetACLBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating login asset ACL: %s", resp.Status)
	}

	return resourceLoginAssetACLRead(ctx, d, m)
}

func resourceLoginAssetACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", loginAssetACLBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting login asset ACL: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
