package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const assetPermUserGroupRelationBasePath = "perms/asset-permissions-user-groups-relations/"

func resourceAssetPermUserGroupRelation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetPermUserGroupRelationCreate,
		ReadContext:   resourceAssetPermUserGroupRelationRead,
		UpdateContext: resourceAssetPermUserGroupRelationUpdate,
		DeleteContext: resourceAssetPermUserGroupRelationDelete,

		Schema: map[string]*schema.Schema{
			"usergroup": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the user group.",
			},
			"assetpermission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the asset permission.",
			},
			"usergroup_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user group.",
			},
			"assetpermission_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the asset permission.",
			},
		},
	}
}

func resourceAssetPermUserGroupRelationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"usergroup":       d.Get("usergroup").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("POST", assetPermUserGroupRelationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset permission user group relation: %s", resp.Status)
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
		return diag.Errorf("Error retrieving relation ID after creation, response: %v", result)
	}

	return resourceAssetPermUserGroupRelationRead(ctx, d, m)
}

func resourceAssetPermUserGroupRelationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", assetPermUserGroupRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading asset permission user group relation: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "usergroup")
	setStringField(d, result, "assetpermission")
	setStringField(d, result, "usergroup_display")
	setStringField(d, result, "assetpermission_display")

	return diags
}

func resourceAssetPermUserGroupRelationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"usergroup":       d.Get("usergroup").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("PUT", assetPermUserGroupRelationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset permission user group relation: %s", resp.Status)
	}

	return resourceAssetPermUserGroupRelationRead(ctx, d, m)
}

func resourceAssetPermUserGroupRelationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", assetPermUserGroupRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset permission user group relation: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
