package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const orgRoleBindingBasePath = "rbac/org-role-bindings/"

func resourceOrgRoleBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgRoleBindingCreate,
		ReadContext:   resourceOrgRoleBindingRead,
		UpdateContext: resourceOrgRoleBindingUpdate,
		DeleteContext: resourceOrgRoleBindingDelete,

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
				Description: "The UUID of the organization role to bind.",
			},
			"org": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the organization.",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scope of the role binding.",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the organization.",
			},
		},
	}
}

func resourceOrgRoleBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user": d.Get("user").(string),
		"role": d.Get("role").(string),
		"org":  d.Get("org").(string),
	}

	resp, err := c.doRequest("POST", orgRoleBindingBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating org role binding: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving org role binding ID after creation, response: %v", result)
	}

	return resourceOrgRoleBindingRead(ctx, d, m)
}

func resourceOrgRoleBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", orgRoleBindingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading org role binding: %s", resp.Status)
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

func resourceOrgRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user": d.Get("user").(string),
		"role": d.Get("role").(string),
		"org":  d.Get("org").(string),
	}

	resp, err := c.doRequest("PUT", orgRoleBindingBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating org role binding: %s", resp.Status)
	}

	return resourceOrgRoleBindingRead(ctx, d, m)
}

func resourceOrgRoleBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", orgRoleBindingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting org role binding: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
