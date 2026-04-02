package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const directoryBasePath = "assets/directories/"

func resourceDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDirectoryCreate,
		ReadContext:   resourceDirectoryRead,
		UpdateContext: resourceDirectoryUpdate,
		DeleteContext: resourceDirectoryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Required: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func resourceDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	zoneID, err := findZoneIDByName(c, d.Get("zone_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	nodeID, err := findNodeIDByName(c, d.Get("node_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"address":   d.Get("address").(string),
		"platform":  d.Get("platform").(int),
		"zone":      zoneID,
		"nodes":     []string{nodeID},
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

	resp, err := c.doRequest("POST", directoryBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_directory", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error creating directory asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving directory asset ID after creation, response: %v", result)
	}

	d.Set("zone_id", zoneID)
	d.Set("node_ids", []string{nodeID})

	return resourceDirectoryRead(ctx, d, m)
}

func resourceDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", directoryBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading directory asset: %s", resp.Status)
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

	return diags
}

func resourceDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	zoneID := d.Get("zone_id").(string)
	nodeIDsRaw := d.Get("node_ids").([]interface{})
	var nodeID string
	if len(nodeIDsRaw) > 0 {
		nodeID = nodeIDsRaw[0].(string)
	}

	if d.HasChange("zone_name") {
		foundID, err := findZoneIDByName(c, d.Get("zone_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		zoneID = foundID
		d.Set("zone_id", foundID)
	}
	if d.HasChange("node_name") {
		foundID, err := findNodeIDByName(c, d.Get("node_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		nodeID = foundID
		d.Set("node_ids", []string{foundID})
	}

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"address":   d.Get("address").(string),
		"platform":  d.Get("platform").(int),
		"zone":      zoneID,
		"nodes":     []string{nodeID},
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

	resp, err := c.doRequest("PUT", directoryBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating directory asset: %s", resp.Status)
	}

	return resourceDirectoryRead(ctx, d, m)
}

func resourceDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", directoryBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting directory asset: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
