package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const systemRoleBasePath = "rbac/system-roles/"

func resourceSystemRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemRoleCreate,
		ReadContext:   resourceSystemRoleRead,
		UpdateContext: resourceSystemRoleUpdate,
		DeleteContext: resourceSystemRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the system role.",
			},
			"permissions": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of permission IDs to assign to this role.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the role.",
			},
			"builtin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a built-in role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the role.",
			},
			"users_amount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of users assigned to this role.",
			},
		},
	}
}

func systemRoleBuildPayload(d *schema.ResourceData) map[string]interface{} {
	data := map[string]interface{}{
		"name":    d.Get("name").(string),
		"comment": d.Get("comment").(string),
	}

	if v, ok := d.GetOk("permissions"); ok {
		perms := v.([]interface{})
		intPerms := make([]int, len(perms))
		for i, p := range perms {
			intPerms[i] = p.(int)
		}
		data["permissions"] = intPerms
	}

	return data
}

func resourceSystemRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := systemRoleBuildPayload(d)

	resp, err := c.doRequest("POST", systemRoleBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating system role: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving system role ID after creation, response: %v", result)
	}

	return resourceSystemRoleRead(ctx, d, m)
}

func resourceSystemRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", systemRoleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading system role: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "comment")
	setStringField(d, result, "display_name")
	setBoolField(d, result, "builtin")
	setIntField(d, result, "users_amount")

	return diags
}

func resourceSystemRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := systemRoleBuildPayload(d)

	resp, err := c.doRequest("PUT", systemRoleBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating system role: %s", resp.Status)
	}

	return resourceSystemRoleRead(ctx, d, m)
}

func resourceSystemRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", systemRoleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting system role: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
