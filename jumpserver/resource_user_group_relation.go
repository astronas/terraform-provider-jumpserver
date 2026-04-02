package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const userGroupRelationBasePath = "users/users-groups-relations/"

func resourceUserGroupRelation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupRelationCreate,
		ReadContext:   resourceUserGroupRelationRead,
		UpdateContext: resourceUserGroupRelationUpdate,
		DeleteContext: resourceUserGroupRelationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the user.",
			},
			"usergroup": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the user group.",
			},
			"user_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user.",
			},
			"usergroup_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user group.",
			},
		},
	}
}

func resourceUserGroupRelationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":      d.Get("user").(string),
		"usergroup": d.Get("usergroup").(string),
	}

	resp, err := c.doRequest("POST", userGroupRelationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating user group relation: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(float64); ok {
		d.SetId(fmt.Sprintf("%d", int(id)))
	} else if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving user group relation ID after creation, response: %v", result)
	}

	return resourceUserGroupRelationRead(ctx, d, m)
}

func resourceUserGroupRelationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", userGroupRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading user group relation: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "user")
	setStringField(d, result, "usergroup")
	setStringField(d, result, "user_display")
	setStringField(d, result, "usergroup_display")

	return diags
}

func resourceUserGroupRelationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":      d.Get("user").(string),
		"usergroup": d.Get("usergroup").(string),
	}

	resp, err := c.doRequest("PUT", userGroupRelationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating user group relation: %s", resp.Status)
	}

	return resourceUserGroupRelationRead(ctx, d, m)
}

func resourceUserGroupRelationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", userGroupRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting user group relation: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
