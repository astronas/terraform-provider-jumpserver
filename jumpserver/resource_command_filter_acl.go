package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const commandFilterACLBasePath = "acls/command-filter-acls/"

func resourceCommandFilterACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommandFilterACLCreate,
		ReadContext:   resourceCommandFilterACLRead,
		UpdateContext: resourceCommandFilterACLUpdate,
		DeleteContext: resourceCommandFilterACLDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the command filter ACL.",
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
			"command_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of command group IDs.",
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
					"reject", "accept", "review",
				}, false),
				Description: "Action when matched: reject, accept, review.",
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

func resourceCommandFilterACLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	// Parse users JSON
	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}
	// Parse assets JSON
	var assets interface{}
	if err := json.Unmarshal([]byte(d.Get("assets").(string)), &assets); err != nil {
		return diag.Errorf("assets must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"users":     users,
		"assets":    assets,
		"accounts":  d.Get("accounts").([]interface{}),
		"priority":  d.Get("priority").(int),
		"action":    d.Get("action").(string),
		"is_active": d.Get("is_active").(bool),
	}

	if v, ok := d.GetOk("command_groups"); ok {
		data["command_groups"] = v.([]interface{})
	}
	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", commandFilterACLBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating command filter ACL: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving command filter ACL ID after creation, response: %v", result)
	}

	return resourceCommandFilterACLRead(ctx, d, m)
}

func resourceCommandFilterACLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", commandFilterACLBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading command filter ACL: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setIntField(d, result, "priority")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")

	// Serialize users/assets back to JSON string
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

	setObjectIDsField(d, result, "command_groups")
	setEnumField(d, result, "action")
	setObjectIDsField(d, result, "reviewers")

	return diags
}

func resourceCommandFilterACLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}
	var assets interface{}
	if err := json.Unmarshal([]byte(d.Get("assets").(string)), &assets); err != nil {
		return diag.Errorf("assets must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"users":     users,
		"assets":    assets,
		"accounts":  d.Get("accounts").([]interface{}),
		"priority":  d.Get("priority").(int),
		"action":    d.Get("action").(string),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("command_groups"); ok {
		data["command_groups"] = v.([]interface{})
	}
	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}

	resp, err := c.doRequest("PUT", commandFilterACLBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating command filter ACL: %s", resp.Status)
	}

	return resourceCommandFilterACLRead(ctx, d, m)
}

func resourceCommandFilterACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", commandFilterACLBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting command filter ACL: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
