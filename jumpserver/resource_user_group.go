package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const userGroupBasePath = "users/groups/"

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user group.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A comment or description for the user group.",
			},
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs to include in this group.",
			},
		},
	}
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	groupData := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("comment"); ok {
		groupData["comment"] = v.(string)
	}

	if v, ok := d.GetOk("users"); ok {
		groupData["users"] = v.([]interface{})
	}

	resp, err := c.doRequest("POST", userGroupBasePath, groupData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_user_group", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating user group: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving user group ID after creation, response: %v", result)
	}

	return resourceUserGroupRead(ctx, d, m)
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", userGroupBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading user group: %s", resp.Status)
	}

	var group map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, group, "name")
	setStringField(d, group, "comment")
	if users, ok := group["users"].([]interface{}); ok {
		d.Set("users", users)
	}

	return diags
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	groupData := map[string]interface{}{
		"name":    d.Get("name").(string),
		"comment": d.Get("comment").(string),
	}

	if v, ok := d.GetOk("users"); ok {
		groupData["users"] = v.([]interface{})
	}

	resp, err := c.doRequest("PUT", userGroupBasePath+d.Id()+"/", groupData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating user group: %s", resp.Status)
	}

	return resourceUserGroupRead(ctx, d, m)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", userGroupBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting user group: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
