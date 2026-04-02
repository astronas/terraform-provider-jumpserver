package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const assetPermissionBasePath = "perms/asset-permissions/"

func resourceAssetPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetPermissionCreate,
		ReadContext:   resourceAssetPermissionRead,
		UpdateContext: resourceAssetPermissionUpdate,
		DeleteContext: resourceAssetPermissionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the asset permission rule.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the permission rule is active.",
			},
			"users": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of user IDs to grant permission to.",
			},
			"user_groups": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of user group IDs to grant permission to.",
			},
			"assets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of asset IDs covered by this permission.",
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of node IDs covered by this permission.",
			},
			"system_users": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of system user IDs allowed by this permission.",
			},
			"actions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of allowed actions (e.g. connect, upload, download).",
			},
			"accounts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of accounts allowed by this permission (e.g. @ALL, @SPEC, or specific account IDs). Replaces system_users in JumpServer v4.",
			},
		},
	}
}

func resourceAssetPermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	permission := map[string]interface{}{
		"name":         d.Get("name").(string),
		"is_active":    d.Get("is_active").(bool),
		"users":        d.Get("users").([]interface{}),
		"user_groups":  d.Get("user_groups").([]interface{}),
		"assets":       d.Get("assets").([]interface{}),
		"nodes":        d.Get("nodes").([]interface{}),
		"system_users": d.Get("system_users").([]interface{}),
		"actions":      d.Get("actions").([]interface{}),
		"accounts":     d.Get("accounts").([]interface{}),
	}

	resp, err := c.doRequestV2("POST", assetPermissionBasePath, permission)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_asset_permission", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset permission: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving asset permission ID after creation, response: %v", result)
	}

	return resourceAssetPermissionRead(ctx, d, m)
}

func resourceAssetPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequestV2("GET", assetPermissionBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error fetching asset permission: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setBoolField(d, result, "is_active")
	if users, ok := result["users"].([]interface{}); ok {
		d.Set("users", users)
	}
	if userGroups, ok := result["user_groups"].([]interface{}); ok {
		d.Set("user_groups", userGroups)
	}
	if assets, ok := result["assets"].([]interface{}); ok {
		d.Set("assets", assets)
	}
	if nodes, ok := result["nodes"].([]interface{}); ok {
		d.Set("nodes", nodes)
	}
	if systemUsers, ok := result["system_users"].([]interface{}); ok {
		d.Set("system_users", systemUsers)
	}
	if actions, ok := result["actions"].([]interface{}); ok {
		d.Set("actions", actions)
	}
	if accounts, ok := result["accounts"].([]interface{}); ok {
		d.Set("accounts", accounts)
	}

	return diags
}

func resourceAssetPermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	permission := map[string]interface{}{
		"name":         d.Get("name").(string),
		"is_active":    d.Get("is_active").(bool),
		"users":        d.Get("users").([]interface{}),
		"user_groups":  d.Get("user_groups").([]interface{}),
		"assets":       d.Get("assets").([]interface{}),
		"nodes":        d.Get("nodes").([]interface{}),
		"system_users": d.Get("system_users").([]interface{}),
		"actions":      d.Get("actions").([]interface{}),
		"accounts":     d.Get("accounts").([]interface{}),
	}

	resp, err := c.doRequestV2("PUT", assetPermissionBasePath+d.Id()+"/", permission)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset permission: %s", resp.Status)
	}

	return resourceAssetPermissionRead(ctx, d, m)
}

func resourceAssetPermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequestV2("DELETE", assetPermissionBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset permission: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
