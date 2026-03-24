package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const assetPermUserRelationBasePath = "perms/asset-permissions-users-relations/"

func resourceAssetPermUserRelation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetPermUserRelationCreate,
		ReadContext:   resourceAssetPermUserRelationRead,
		UpdateContext: resourceAssetPermUserRelationUpdate,
		DeleteContext: resourceAssetPermUserRelationDelete,

		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the user.",
			},
			"assetpermission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the asset permission.",
			},
			"user_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user.",
			},
			"assetpermission_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the asset permission.",
			},
		},
	}
}

func resourceAssetPermUserRelationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":            d.Get("user").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("POST", assetPermUserRelationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset permission user relation: %s", resp.Status)
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

	return resourceAssetPermUserRelationRead(ctx, d, m)
}

func resourceAssetPermUserRelationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", assetPermUserRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading asset permission user relation: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "user")
	setStringField(d, result, "assetpermission")
	setStringField(d, result, "user_display")
	setStringField(d, result, "assetpermission_display")

	return diags
}

func resourceAssetPermUserRelationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":            d.Get("user").(string),
		"assetpermission": d.Get("assetpermission").(string),
	}

	resp, err := c.doRequest("PUT", assetPermUserRelationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset permission user relation: %s", resp.Status)
	}

	return resourceAssetPermUserRelationRead(ctx, d, m)
}

func resourceAssetPermUserRelationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", assetPermUserRelationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset permission user relation: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
