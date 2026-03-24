package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const connectMethodACLBasePath = "acls/connect-method-acls/"

func resourceConnectMethodACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectMethodACLCreate,
		ReadContext:   resourceConnectMethodACLRead,
		UpdateContext: resourceConnectMethodACLUpdate,
		DeleteContext: resourceConnectMethodACLDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connect method ACL.",
			},
			"users": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User filter as JSON (e.g. {\"type\":\"ids\",\"ids\":[\"uuid1\"]}).",
			},
			"connect_methods": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of allowed connect method identifiers.",
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

func resourceConnectMethodACLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"users":     users,
		"priority":  d.Get("priority").(int),
		"action":    d.Get("action").(string),
		"is_active": d.Get("is_active").(bool),
	}

	if v, ok := d.GetOk("connect_methods"); ok {
		data["connect_methods"] = v.([]interface{})
	}
	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", connectMethodACLBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating connect method ACL: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving connect method ACL ID after creation, response: %v", result)
	}

	return resourceConnectMethodACLRead(ctx, d, m)
}

func resourceConnectMethodACLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", connectMethodACLBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading connect method ACL: %s", resp.Status)
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

	if v, ok := result["connect_methods"].([]interface{}); ok {
		d.Set("connect_methods", v)
	}

	setEnumField(d, result, "action")
	setObjectIDsField(d, result, "reviewers")

	return diags
}

func resourceConnectMethodACLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"users":     users,
		"priority":  d.Get("priority").(int),
		"action":    d.Get("action").(string),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("connect_methods"); ok {
		data["connect_methods"] = v.([]interface{})
	}
	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}

	resp, err := c.doRequest("PUT", connectMethodACLBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating connect method ACL: %s", resp.Status)
	}

	return resourceConnectMethodACLRead(ctx, d, m)
}

func resourceConnectMethodACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", connectMethodACLBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting connect method ACL: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
