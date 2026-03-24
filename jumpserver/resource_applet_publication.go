package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const appletPublicationBasePath = "terminal/applet-publications/"

func resourceAppletPublication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppletPublicationCreate,
		ReadContext:   resourceAppletPublicationRead,
		UpdateContext: resourceAppletPublicationUpdate,
		DeleteContext: resourceAppletPublicationDelete,

		Schema: map[string]*schema.Schema{
			"applet": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the applet.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the applet host.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "pending",
				Description: "The publication status.",
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

func resourceAppletPublicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"applet": d.Get("applet").(string),
		"host":   d.Get("host").(string),
		"status": d.Get("status").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", appletPublicationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating applet publication: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving applet publication ID after creation, response: %v", result)
	}

	return resourceAppletPublicationRead(ctx, d, m)
}

func resourceAppletPublicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", appletPublicationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading applet publication: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "applet")
	setObjectIDField(d, result, "host")
	setEnumField(d, result, "status")
	setStringField(d, result, "comment")

	return diags
}

func resourceAppletPublicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"applet":  d.Get("applet").(string),
		"host":    d.Get("host").(string),
		"status":  d.Get("status").(string),
		"comment": d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", appletPublicationBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating applet publication: %s", resp.Status)
	}

	return resourceAppletPublicationRead(ctx, d, m)
}

func resourceAppletPublicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", appletPublicationBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting applet publication: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
