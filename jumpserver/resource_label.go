package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const labelBasePath = "labels/labels/"

func resourceLabel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLabelCreate,
		ReadContext:   resourceLabelRead,
		UpdateContext: resourceLabelUpdate,
		DeleteContext: resourceLabelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name (key) of the label.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the label.",
			},
			"color": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The display color of the label.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the label.",
			},
		},
	}
}

func resourceLabelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	labelData := map[string]interface{}{
		"name":  d.Get("name").(string),
		"value": d.Get("value").(string),
	}

	if v, ok := d.GetOk("color"); ok {
		labelData["color"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		labelData["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", labelBasePath, labelData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating label: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving label ID after creation, response: %v", result)
	}

	return resourceLabelRead(ctx, d, m)
}

func resourceLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", labelBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading label: %s", resp.Status)
	}

	var label map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&label); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, label, "name")
	setStringField(d, label, "value")
	setStringField(d, label, "color")
	setStringField(d, label, "comment")

	return diags
}

func resourceLabelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	labelData := map[string]interface{}{
		"name":    d.Get("name").(string),
		"value":   d.Get("value").(string),
		"color":   d.Get("color").(string),
		"comment": d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", labelBasePath+d.Id()+"/", labelData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating label: %s", resp.Status)
	}

	return resourceLabelRead(ctx, d, m)
}

func resourceLabelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", labelBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting label: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
