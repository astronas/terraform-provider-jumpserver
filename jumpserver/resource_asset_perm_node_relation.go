package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const assetPermNodeRelationBasePath = "perms/asset-permissions-nodes-relations/"

func resourceAssetPermNodeRelation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetPermNodeRelationCreate,
		ReadContext:   resourceAssetPermNodeRelationRead,
		UpdateContext: resourceAssetPermNodeRelationUpdate,
		DeleteContext: resourceAssetPermNodeRelationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"node": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the node.",
			},
			"assetpermission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the asset permission.",
			},
			"node_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the node.",
			},
			"assetpermission_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the asset permission.",
			},
		},
	}
}

func resourceAssetPermNodeRelationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"node":            d.Get("node").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("POST", assetPermNodeRelationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset permission node relation: %s", resp.Status)
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

	return resourceAssetPermNodeRelationRead(ctx, d, m)
}

func resourceAssetPermNodeRelationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", assetPermNodeRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading asset permission node relation: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "node")
	setStringField(d, result, "assetpermission")
	setStringField(d, result, "node_display")
	setStringField(d, result, "assetpermission_display")

	return diags
}

func resourceAssetPermNodeRelationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"node":            d.Get("node").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("PUT", assetPermNodeRelationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset permission node relation: %s", resp.Status)
	}

	return resourceAssetPermNodeRelationRead(ctx, d, m)
}

func resourceAssetPermNodeRelationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", assetPermNodeRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset permission node relation: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
