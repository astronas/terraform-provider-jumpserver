package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const virtualAppPublicationBasePath = "terminal/virtual-app-publications/"

func resourceVirtualAppPublication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualAppPublicationCreate,
		ReadContext:   resourceVirtualAppPublicationRead,
		UpdateContext: resourceVirtualAppPublicationUpdate,
		DeleteContext: resourceVirtualAppPublicationDelete,

		Schema: map[string]*schema.Schema{
			"app": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the virtual app.",
			},
			"provider_host": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the app-provider host.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "pending",
				Description: "Status of the publication.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the publication.",
			},
		},
	}
}

func resourceVirtualAppPublicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"app":      d.Get("app").(string),
		"provider": d.Get("provider_host").(string),
		"status":   d.Get("status").(string),
		"comment":  d.Get("comment").(string),
	}

	resp, err := c.doRequest("POST", virtualAppPublicationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating virtual app publication: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving virtual app publication ID after creation, response: %v", result)
	}

	return resourceVirtualAppPublicationRead(ctx, d, m)
}

func resourceVirtualAppPublicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", virtualAppPublicationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading virtual app publication: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "app")
	if v, ok := readObjectID(result, "provider"); ok {
		d.Set("provider_host", v)
	}
	setStringField(d, result, "status")
	setStringField(d, result, "comment")

	return diags
}

func resourceVirtualAppPublicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"app":      d.Get("app").(string),
		"provider": d.Get("provider_host").(string),
		"status":   d.Get("status").(string),
		"comment":  d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", virtualAppPublicationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating virtual app publication: %s", resp.Status)
	}

	return resourceVirtualAppPublicationRead(ctx, d, m)
}

func resourceVirtualAppPublicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", virtualAppPublicationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting virtual app publication: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
