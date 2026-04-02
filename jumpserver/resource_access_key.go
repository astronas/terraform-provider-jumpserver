package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const accessKeyBasePath = "authentication/access-keys/"

func resourceAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessKeyCreate,
		ReadContext:   resourceAccessKeyRead,
		UpdateContext: resourceAccessKeyUpdate,
		DeleteContext: resourceAccessKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the access key is active.",
			},
			"ip_group": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of allowed IP addresses or ranges. Default is all ('*').",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date of the access key.",
			},
			"date_last_used": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last used date of the access key.",
			},
		},
	}
}

func resourceAccessKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"is_active": d.Get("is_active").(bool),
	}
	if v, ok := d.GetOk("ip_group"); ok {
		data["ip_group"] = v.([]interface{})
	} else {
		data["ip_group"] = []string{"*"}
	}

	resp, err := c.doRequest("POST", accessKeyBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating access key: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving access key ID after creation, response: %v", result)
	}

	return resourceAccessKeyRead(ctx, d, m)
}

func resourceAccessKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", accessKeyBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading access key: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setBoolField(d, result, "is_active")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_last_used")
	if v, ok := result["ip_group"].([]interface{}); ok {
		d.Set("ip_group", v)
	}

	return diags
}

func resourceAccessKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"is_active": d.Get("is_active").(bool),
	}
	if v, ok := d.GetOk("ip_group"); ok {
		data["ip_group"] = v.([]interface{})
	}

	resp, err := c.doRequest("PUT", accessKeyBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating access key: %s", resp.Status)
	}

	return resourceAccessKeyRead(ctx, d, m)
}

func resourceAccessKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", accessKeyBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting access key: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
