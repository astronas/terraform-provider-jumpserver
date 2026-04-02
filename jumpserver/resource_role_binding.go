package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const roleBindingBasePath = "rbac/role-bindings/"

func resourceRoleBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleBindingCreate,
		ReadContext:   resourceRoleBindingRead,
		UpdateContext: resourceRoleBindingUpdate,
		DeleteContext: resourceRoleBindingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the user to bind the role to.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the role to bind.",
			},
			"org": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the organization. Required for org-scoped role bindings.",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scope of the role binding (system or org).",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the organization.",
			},
		},
	}
}

func resourceRoleBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user": d.Get("user").(string),
		"role": d.Get("role").(string),
	}

	if v, ok := d.GetOk("org"); ok {
		data["org"] = v.(string)
	}

	resp, err := c.doRequest("POST", roleBindingBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating role binding: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving role binding ID after creation, response: %v", result)
	}

	return resourceRoleBindingRead(ctx, d, m)
}

func resourceRoleBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", roleBindingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading role binding: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "user")
	setStringField(d, result, "role")
	setStringField(d, result, "org")
	setEnumField(d, result, "scope")
	setStringField(d, result, "org_name")

	return diags
}

func resourceRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user": d.Get("user").(string),
		"role": d.Get("role").(string),
	}

	if v, ok := d.GetOk("org"); ok {
		data["org"] = v.(string)
	}

	resp, err := c.doRequest("PUT", roleBindingBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating role binding: %s", resp.Status)
	}

	return resourceRoleBindingRead(ctx, d, m)
}

func resourceRoleBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", roleBindingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting role binding: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
