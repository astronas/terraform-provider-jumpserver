package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const favoriteAssetBasePath = "assets/favorite-assets/"

func resourceFavoriteAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFavoriteAssetCreate,
		ReadContext:   resourceFavoriteAssetRead,
		DeleteContext: resourceFavoriteAssetDelete,

		Schema: map[string]*schema.Schema{
			"asset": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the asset to favorite.",
			},
		},
	}
}

func resourceFavoriteAssetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"asset": d.Get("asset").(string),
	}

	resp, err := c.doRequest("POST", favoriteAssetBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating favorite asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else if asset, ok := result["asset"].(string); ok {
		d.SetId(asset)
	} else {
		return diag.Errorf("Error retrieving favorite asset ID after creation, response: %v", result)
	}

	return nil
}

func resourceFavoriteAssetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", favoriteAssetBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading favorite asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "asset")

	return diags
}

func resourceFavoriteAssetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", favoriteAssetBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting favorite asset: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
