package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const zoneBasePath = "assets/zones/"

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the zone.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the zone.",
			},
		},
	}
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	zoneData := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("comment"); ok {
		zoneData["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", zoneBasePath, zoneData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_zone", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating zone: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving zone ID after creation, response: %v", result)
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", zoneBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading zone: %s", resp.Status)
	}

	var zone map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&zone); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, zone, "name")
	setStringField(d, zone, "comment")

	return diags
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	zoneData := map[string]interface{}{
		"name":    d.Get("name").(string),
		"comment": d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", zoneBasePath+d.Id()+"/", zoneData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating zone: %s", resp.Status)
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", zoneBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting zone: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
