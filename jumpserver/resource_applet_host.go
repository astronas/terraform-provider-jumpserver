package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const appletHostBasePath = "terminal/applet-hosts/"

func resourceAppletHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppletHostCreate,
		ReadContext:   resourceAppletHostRead,
		UpdateContext: resourceAppletHostUpdate,
		DeleteContext: resourceAppletHostDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_invalid": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "error",
						},
						"is_active": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"protocols": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"deploy_options": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of deployment options.",
			},
		},
	}
}

func appletHostBuildPayload(d *schema.ResourceData) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"address":   d.Get("address").(string),
		"platform":  d.Get("platform").(int),
		"is_active": d.Get("is_active").(bool),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}
	if v, ok := d.GetOk("accounts"); ok {
		data["accounts"] = expandAccounts(v.([]interface{}))
	}
	if v, ok := d.GetOk("protocols"); ok {
		data["protocols"] = expandProtocols(v.([]interface{}))
	}
	if v, ok := d.GetOk("deploy_options"); ok {
		var opts interface{}
		if err := json.Unmarshal([]byte(v.(string)), &opts); err != nil {
			return nil, fmt.Errorf("invalid deploy_options JSON: %w", err)
		}
		data["deploy_options"] = opts
	}
	return data, nil
}

func appletHostResolveZoneNode(c *Config, d *schema.ResourceData, data map[string]interface{}) error {
	if v, ok := d.GetOk("zone_name"); ok {
		zoneID, err := findZoneIDByName(c, v.(string))
		if err != nil {
			return err
		}
		data["zone"] = zoneID
		d.Set("zone_id", zoneID)
	}
	if v, ok := d.GetOk("node_name"); ok {
		nodeID, err := findNodeIDByName(c, v.(string))
		if err != nil {
			return err
		}
		data["nodes"] = []string{nodeID}
		d.Set("node_ids", []string{nodeID})
	}
	return nil
}

func appletHostUpdateZoneNode(c *Config, d *schema.ResourceData, data map[string]interface{}) error {
	if d.HasChange("zone_name") {
		if v, ok := d.GetOk("zone_name"); ok {
			foundID, err := findZoneIDByName(c, v.(string))
			if err != nil {
				return err
			}
			data["zone"] = foundID
			d.Set("zone_id", foundID)
		}
	} else if v := d.Get("zone_id").(string); v != "" {
		data["zone"] = v
	}
	if d.HasChange("node_name") {
		if v, ok := d.GetOk("node_name"); ok {
			foundID, err := findNodeIDByName(c, v.(string))
			if err != nil {
				return err
			}
			data["nodes"] = []string{foundID}
			d.Set("node_ids", []string{foundID})
		}
	} else if v, ok := d.GetOk("node_ids"); ok {
		data["nodes"] = v.([]interface{})
	}
	return nil
}

func resourceAppletHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data, err := appletHostBuildPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := appletHostResolveZoneNode(c, d, data); err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.doRequest("POST", appletHostBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error creating applet host: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving applet host ID after creation, response: %v", result)
	}

	return resourceAppletHostRead(ctx, d, m)
}

func resourceAppletHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", appletHostBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading applet host: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "address")
	setStringField(d, result, "comment")
	setIntField(d, result, "platform")
	setBoolField(d, result, "is_active")
	if v, ok := result["zone"].(string); ok {
		d.Set("zone_id", v)
	}
	if v, ok := result["nodes"].([]interface{}); ok {
		d.Set("node_ids", v)
	}
	if v, ok := result["accounts"].([]interface{}); ok {
		d.Set("accounts", flattenAccounts(v))
	}
	if v, ok := result["protocols"].([]interface{}); ok {
		d.Set("protocols", flattenProtocols(v))
	}
	if v, ok := result["deploy_options"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("deploy_options", string(jsonBytes))
	}

	return diags
}

func resourceAppletHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data, err := appletHostBuildPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := appletHostUpdateZoneNode(c, d, data); err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.doRequest("PUT", appletHostBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating applet host: %s", resp.Status)
	}

	return resourceAppletHostRead(ctx, d, m)
}

func resourceAppletHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", appletHostBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting applet host: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
