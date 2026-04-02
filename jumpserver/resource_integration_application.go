package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const integrationAppBasePath = "accounts/integration-applications/"

func resourceIntegrationApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationApplicationCreate,
		ReadContext:   resourceIntegrationApplicationRead,
		UpdateContext: resourceIntegrationApplicationUpdate,
		DeleteContext: resourceIntegrationApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the integration application.",
			},
			"accounts": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of account IDs associated with this application.",
			},
			"ip_group": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of allowed IP ranges (e.g. [\"*\"] for all). Defaults to [\"*\"].",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the application is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the application.",
			},
		},
	}
}

func resourceIntegrationApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"accounts":  d.Get("accounts").([]interface{}),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("ip_group"); ok {
		data["ip_group"] = v.([]interface{})
	} else {
		data["ip_group"] = []string{"*"}
	}

	resp, err := c.doRequest("POST", integrationAppBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_integration_application", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating integration application: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving integration application ID after creation, response: %v", result)
	}

	return resourceIntegrationApplicationRead(ctx, d, m)
}

func resourceIntegrationApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", integrationAppBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading integration application: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")
	setObjectIDsField(d, result, "accounts")

	// ip_group: array of strings
	if v, ok := result["ip_group"].([]interface{}); ok {
		d.Set("ip_group", v)
	}

	return diags
}

func resourceIntegrationApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"accounts":  d.Get("accounts").([]interface{}),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("ip_group"); ok {
		data["ip_group"] = v.([]interface{})
	} else {
		data["ip_group"] = []string{"*"}
	}

	resp, err := c.doRequest("PUT", integrationAppBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating integration application: %s", resp.Status)
	}

	return resourceIntegrationApplicationRead(ctx, d, m)
}

func resourceIntegrationApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", integrationAppBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting integration application: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
