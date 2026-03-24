package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const appletHostDeploymentBasePath = "terminal/applet-host-deployments/"

func resourceAppletHostDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppletHostDeploymentCreate,
		ReadContext:   resourceAppletHostDeploymentRead,
		UpdateContext: resourceAppletHostDeploymentUpdate,
		DeleteContext: resourceAppletHostDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the applet host to deploy to.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the deployment.",
			},
			"install_applets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to install applets during deployment.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
			"date_start": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment start date.",
			},
			"date_finished": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment finish date.",
			},
		},
	}
}

func resourceAppletHostDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"host":            d.Get("host").(string),
		"install_applets": d.Get("install_applets").(bool),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", appletHostDeploymentBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating applet host deployment: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving applet host deployment ID after creation, response: %v", result)
	}

	return resourceAppletHostDeploymentRead(ctx, d, m)
}

func resourceAppletHostDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", appletHostDeploymentBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading applet host deployment: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "host")
	setStringField(d, result, "comment")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_start")
	setStringField(d, result, "date_finished")

	return diags
}

func resourceAppletHostDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"host": d.Get("host").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", appletHostDeploymentBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating applet host deployment: %s", resp.Status)
	}

	return resourceAppletHostDeploymentRead(ctx, d, m)
}

func resourceAppletHostDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", appletHostDeploymentBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting applet host deployment: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
