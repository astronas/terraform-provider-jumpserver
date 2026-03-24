package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const assetPermAssetRelationBasePath = "perms/asset-permissions-assets-relations/"

func resourceAssetPermAssetRelation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetPermAssetRelationCreate,
		ReadContext:   resourceAssetPermAssetRelationRead,
		UpdateContext: resourceAssetPermAssetRelationUpdate,
		DeleteContext: resourceAssetPermAssetRelationDelete,

		Schema: map[string]*schema.Schema{
			"asset": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the asset.",
			},
			"assetpermission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the asset permission.",
			},
			"asset_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the asset.",
			},
			"assetpermission_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the asset permission.",
			},
		},
	}
}

func resourceAssetPermAssetRelationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"asset":           d.Get("asset").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("POST", assetPermAssetRelationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset permission asset relation: %s", resp.Status)
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

	return resourceAssetPermAssetRelationRead(ctx, d, m)
}

func resourceAssetPermAssetRelationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", assetPermAssetRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading asset permission asset relation: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "asset")
	setStringField(d, result, "assetpermission")
	setStringField(d, result, "asset_display")
	setStringField(d, result, "assetpermission_display")

	return diags
}

func resourceAssetPermAssetRelationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"asset":           d.Get("asset").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("PUT", assetPermAssetRelationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset permission asset relation: %s", resp.Status)
	}

	return resourceAssetPermAssetRelationRead(ctx, d, m)
}

func resourceAssetPermAssetRelationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", assetPermAssetRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset permission asset relation: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
