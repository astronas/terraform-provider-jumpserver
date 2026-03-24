package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const assetBasePath = "assets/assets/"

func resourceAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetCreate,
		ReadContext:   resourceAssetRead,
		UpdateContext: resourceAssetUpdate,
		DeleteContext: resourceAssetDelete,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname of the asset.",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address of the asset.",
			},
			"platform": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The platform identifier for the asset (e.g. 'Linux', 'Windows').",
			},
			"protocols": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of protocols supported by the asset.",
			},
			"nodes_display": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of node display paths the asset belongs to.",
			},
		},
	}
}

func resourceAssetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	asset := map[string]interface{}{
		"hostname":      d.Get("hostname").(string),
		"ip":            d.Get("ip").(string),
		"platform":      d.Get("platform").(string),
		"protocols":     d.Get("protocols").([]interface{}),
		"nodes_display": d.Get("nodes_display").([]interface{}),
	}

	resp, err := c.doRequest("POST", assetBasePath, asset)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving asset ID after creation, response: %v", result)
	}

	return resourceAssetRead(ctx, d, m)
}

func resourceAssetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", assetBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error fetching asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "hostname")
	setStringField(d, result, "ip")
	setStringField(d, result, "platform")
	if protocols, ok := result["protocols"].([]interface{}); ok {
		d.Set("protocols", protocols)
	}
	if nodesDisplay, ok := result["nodes_display"].([]interface{}); ok {
		d.Set("nodes_display", nodesDisplay)
	}

	return diags
}

func resourceAssetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	asset := map[string]interface{}{
		"hostname":      d.Get("hostname").(string),
		"ip":            d.Get("ip").(string),
		"platform":      d.Get("platform").(string),
		"protocols":     d.Get("protocols").([]interface{}),
		"nodes_display": d.Get("nodes_display").([]interface{}),
	}

	resp, err := c.doRequest("PUT", assetBasePath+d.Id()+"/", asset)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset: %s", resp.Status)
	}

	return resourceAssetRead(ctx, d, m)
}

func resourceAssetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", assetBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
