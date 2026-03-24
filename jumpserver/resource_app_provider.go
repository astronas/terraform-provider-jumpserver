package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const appProviderBasePath = "terminal/app-providers/"

func resourceAppProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppProviderCreate,
		ReadContext:   resourceAppProviderRead,
		UpdateContext: resourceAppProviderUpdate,
		DeleteContext: resourceAppProviderDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the app provider.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname of the app provider.",
			},
			"terminal": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the associated terminal.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date.",
			},
		},
	}
}

func resourceAppProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":     d.Get("name").(string),
		"hostname": d.Get("hostname").(string),
	}
	if v, ok := d.GetOk("terminal"); ok {
		data["terminal"] = v.(string)
	}

	resp, err := c.doRequest("POST", appProviderBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating app provider: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving app provider ID after creation, response: %v", result)
	}

	return resourceAppProviderRead(ctx, d, m)
}

func resourceAppProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", appProviderBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading app provider: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "hostname")
	setStringField(d, result, "terminal")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceAppProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":     d.Get("name").(string),
		"hostname": d.Get("hostname").(string),
	}
	if v, ok := d.GetOk("terminal"); ok {
		data["terminal"] = v.(string)
	}

	resp, err := c.doRequest("PUT", appProviderBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating app provider: %s", resp.Status)
	}

	return resourceAppProviderRead(ctx, d, m)
}

func resourceAppProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", appProviderBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting app provider: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
