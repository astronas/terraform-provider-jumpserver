package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const endpointBasePath = "terminal/endpoints/"

func resourceEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointCreate,
		ReadContext:   resourceEndpointRead,
		UpdateContext: resourceEndpointUpdate,
		DeleteContext: resourceEndpointDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the endpoint.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access address. Empty uses the browser address.",
			},
			"https_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      443,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "HTTPS port.",
			},
			"http_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "HTTP port.",
			},
			"ssh_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2222,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "SSH protocol port.",
			},
			"rdp_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3389,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "RDP protocol port.",
			},
			"mysql_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      33061,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "MySQL protocol port.",
			},
			"mariadb_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      33062,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "MariaDB protocol port.",
			},
			"postgresql_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      54320,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "PostgreSQL protocol port.",
			},
			"redis_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      63790,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "Redis protocol port.",
			},
			"vnc_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15900,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "VNC protocol port.",
			},
			"oracle_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15210,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "Oracle protocol port.",
			},
			"sqlserver_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      14330,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "SQL Server protocol port.",
			},
			"mongodb_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      27018,
				ValidateFunc: validation.IntBetween(0, 65535),
				Description:  "MongoDB protocol port.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the endpoint is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the endpoint.",
			},
		},
	}
}

var endpointPortFields = []string{
	"https_port", "http_port", "ssh_port", "rdp_port",
	"mysql_port", "mariadb_port", "postgresql_port", "redis_port",
	"vnc_port", "oracle_port", "sqlserver_port", "mongodb_port",
}

func endpointBuildPayload(d *schema.ResourceData) map[string]interface{} {
	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("host"); ok {
		data["host"] = v.(string)
	}

	for _, field := range endpointPortFields {
		data[field] = d.Get(field).(int)
	}

	return data
}

func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	resp, err := c.doRequest("POST", endpointBasePath, endpointBuildPayload(d))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating endpoint: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving endpoint ID after creation, response: %v", result)
	}

	return resourceEndpointRead(ctx, d, m)
}

func resourceEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", endpointBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading endpoint: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "host")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")

	for _, field := range endpointPortFields {
		setIntField(d, result, field)
	}

	return diags
}

func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	resp, err := c.doRequest("PUT", endpointBasePath+d.Id()+"/", endpointBuildPayload(d))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating endpoint: %s", resp.Status)
	}

	return resourceEndpointRead(ctx, d, m)
}

func resourceEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", endpointBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting endpoint: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
