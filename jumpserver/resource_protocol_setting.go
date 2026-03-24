package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const protocolSettingBasePath = "assets/protocol-settings/"

func resourceProtocolSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtocolSettingCreate,
		ReadContext:   resourceProtocolSettingRead,
		UpdateContext: resourceProtocolSettingUpdate,
		DeleteContext: resourceProtocolSettingDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol name.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The default port number (0-65535).",
			},
			"primary": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this protocol is primary.",
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this protocol is required.",
			},
			"default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this protocol is default.",
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this protocol is public.",
			},
			"setting": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of protocol settings.",
			},
			"port_from_addr": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the port is derived from the address.",
			},
		},
	}
}

func resourceProtocolSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":     d.Get("name").(string),
		"port":     d.Get("port").(int),
		"primary":  d.Get("primary").(bool),
		"required": d.Get("required").(bool),
		"default":  d.Get("default").(bool),
		"public":   d.Get("public").(bool),
	}
	if v, ok := d.GetOk("setting"); ok {
		var setting interface{}
		if err := json.Unmarshal([]byte(v.(string)), &setting); err != nil {
			return diag.Errorf("invalid setting JSON: %s", err)
		}
		data["setting"] = setting
	}

	resp, err := c.doRequest("POST", protocolSettingBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating protocol setting: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(float64); ok {
		d.SetId(fmt.Sprintf("%d", int(id)))
	} else if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving protocol setting ID after creation, response: %v", result)
	}

	return resourceProtocolSettingRead(ctx, d, m)
}

func resourceProtocolSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", protocolSettingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading protocol setting: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setIntField(d, result, "port")
	setBoolField(d, result, "primary")
	setBoolField(d, result, "required")
	setBoolField(d, result, "default")
	setBoolField(d, result, "public")
	setBoolField(d, result, "port_from_addr")
	if v, ok := result["setting"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("setting", string(jsonBytes))
	}

	return diags
}

func resourceProtocolSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":     d.Get("name").(string),
		"port":     d.Get("port").(int),
		"primary":  d.Get("primary").(bool),
		"required": d.Get("required").(bool),
		"default":  d.Get("default").(bool),
		"public":   d.Get("public").(bool),
	}
	if v, ok := d.GetOk("setting"); ok {
		var setting interface{}
		if err := json.Unmarshal([]byte(v.(string)), &setting); err != nil {
			return diag.Errorf("invalid setting JSON: %s", err)
		}
		data["setting"] = setting
	}

	resp, err := c.doRequest("PUT", protocolSettingBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating protocol setting: %s", resp.Status)
	}

	return resourceProtocolSettingRead(ctx, d, m)
}

func resourceProtocolSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", protocolSettingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting protocol setting: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
