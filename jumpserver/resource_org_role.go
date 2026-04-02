package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const orgRoleBasePath = "rbac/org-roles/"

func resourceOrgRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgRoleCreate,
		ReadContext:   resourceOrgRoleRead,
		UpdateContext: resourceOrgRoleUpdate,
		DeleteContext: resourceOrgRoleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization role.",
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

func orgRoleBuildPayload(d *schema.ResourceData) map[string]interface{} {
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

func resourceOrgRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := orgRoleBuildPayload(d)

	resp, err := c.doRequest("POST", orgRoleBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_org_role", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating org role: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving org role ID after creation, response: %v", result)
	}

	return resourceOrgRoleRead(ctx, d, m)
}

func resourceOrgRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", orgRoleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading org role: %s", resp.Status)
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

func resourceOrgRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := orgRoleBuildPayload(d)

	resp, err := c.doRequest("PUT", orgRoleBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating org role: %s", resp.Status)
	}

	return resourceOrgRoleRead(ctx, d, m)
}

func resourceOrgRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", orgRoleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting org role: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
