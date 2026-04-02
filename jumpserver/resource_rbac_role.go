package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const rbacRoleBasePath = "rbac/roles/"

func resourceRbacRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRbacRoleCreate,
		ReadContext:   resourceRbacRoleRead,
		UpdateContext: resourceRbacRoleUpdate,
		DeleteContext: resourceRbacRoleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the RBAC role.",
			},
			"permissions": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of permission IDs to assign to this role (write-only).",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the role.",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scope of the role.",
			},
			"builtin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a built-in role.",
			},
			"users_amount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of users assigned to this role.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Who created this role.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Who last updated this role.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date.",
			},
		},
	}
}

func rbacRoleBuildPayload(d *schema.ResourceData) map[string]interface{} {
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

func resourceRbacRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := rbacRoleBuildPayload(d)

	resp, err := c.doRequest("POST", rbacRoleBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_rbac_role", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating RBAC role: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving RBAC role ID after creation, response: %v", result)
	}

	return resourceRbacRoleRead(ctx, d, m)
}

func resourceRbacRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", rbacRoleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading RBAC role: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "comment")
	setStringField(d, result, "display_name")
	setEnumField(d, result, "scope")
	setBoolField(d, result, "builtin")
	setIntField(d, result, "users_amount")
	setStringField(d, result, "created_by")
	setStringField(d, result, "updated_by")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceRbacRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := rbacRoleBuildPayload(d)

	resp, err := c.doRequest("PUT", rbacRoleBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating RBAC role: %s", resp.Status)
	}

	return resourceRbacRoleRead(ctx, d, m)
}

func resourceRbacRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", rbacRoleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting RBAC role: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
