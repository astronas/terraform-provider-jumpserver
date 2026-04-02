package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const labeledResourceBasePath = "labels/labeled-resources/"

func resourceLabeledResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLabeledResourceCreate,
		ReadContext:   resourceLabeledResourceRead,
		UpdateContext: resourceLabeledResourceUpdate,
		DeleteContext: resourceLabeledResourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the label.",
			},
			"res_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource type (e.g. asset, user).",
			},
			"res_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the resource to label.",
			},
		},
	}
}

func resourceLabeledResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"label":    d.Get("label").(string),
		"res_type": d.Get("res_type").(string),
		"res_id":   d.Get("res_id").(string),
	}

	resp, err := c.doRequest("POST", labeledResourceBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating labeled resource: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else if id, ok := result["id"].(float64); ok {
		d.SetId(fmt.Sprintf("%d", int(id)))
	} else {
		return diag.Errorf("Error retrieving labeled resource ID after creation, response: %v", result)
	}

	return resourceLabeledResourceRead(ctx, d, m)
}

func resourceLabeledResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", labeledResourceBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading labeled resource: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "label")
	setStringField(d, result, "res_type")
	setStringField(d, result, "res_id")

	return diags
}

func resourceLabeledResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"label":    d.Get("label").(string),
		"res_type": d.Get("res_type").(string),
		"res_id":   d.Get("res_id").(string),
	}

	resp, err := c.doRequest("PUT", labeledResourceBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating labeled resource: %s", resp.Status)
	}

	return resourceLabeledResourceRead(ctx, d, m)
}

func resourceLabeledResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", labeledResourceBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting labeled resource: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
