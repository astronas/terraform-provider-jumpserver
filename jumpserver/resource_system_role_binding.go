package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const systemRoleBindingBasePath = "rbac/system-role-bindings/"

func resourceSystemRoleBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemRoleBindingCreate,
		ReadContext:   resourceSystemRoleBindingRead,
		UpdateContext: resourceSystemRoleBindingUpdate,
		DeleteContext: resourceSystemRoleBindingDelete,

		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the user to bind the role to.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the system role to bind.",
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

func resourceSystemRoleBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user": d.Get("user").(string),
		"role": d.Get("role").(string),
	}

	resp, err := c.doRequest("POST", systemRoleBindingBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating system role binding: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving system role binding ID after creation, response: %v", result)
	}

	return resourceSystemRoleBindingRead(ctx, d, m)
}

func resourceSystemRoleBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", systemRoleBindingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading system role binding: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "user")
	setStringField(d, result, "role")
	setEnumField(d, result, "scope")
	setStringField(d, result, "org_name")

	return diags
}

func resourceSystemRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user": d.Get("user").(string),
		"role": d.Get("role").(string),
	}

	resp, err := c.doRequest("PUT", systemRoleBindingBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating system role binding: %s", resp.Status)
	}

	return resourceSystemRoleBindingRead(ctx, d, m)
}

func resourceSystemRoleBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", systemRoleBindingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting system role binding: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
