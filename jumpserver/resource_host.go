package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const hostBasePath = "assets/hosts/"

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
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

			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

// -------------------------------------------------------------------
// Create
// -------------------------------------------------------------------
func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	zoneName := d.Get("zone_name").(string)
	zoneID, err := findZoneIDByName(c, zoneName)
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get("node_name").(string)
	nodeID, err := findNodeIDByName(c, nodeName)
	if err != nil {
		return diag.FromErr(err)
	}

	hostData := map[string]interface{}{
		"name":     d.Get("name").(string),
		"address":  d.Get("address").(string),
		"platform": d.Get("platform").(int),
		"zone":     zoneID,
		"nodes":    []string{nodeID},
	}

	if v, ok := d.GetOk("comment"); ok {
		hostData["comment"] = v.(string)
	}

	if v, ok := d.GetOk("accounts"); ok {
		hostData["accounts"] = expandAccounts(v.([]interface{}))
	}

	if v, ok := d.GetOk("protocols"); ok {
		hostData["protocols"] = expandProtocols(v.([]interface{}))
	}

	resp, err := c.doRequest("POST", hostBasePath, hostData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to create host in JumpServer. HTTP status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	hostID, ok := result["id"].(string)
	if !ok {
		return diag.Errorf("No 'id' field found in host creation response")
	}
	d.SetId(hostID)

	d.Set("zone_id", zoneID)
	d.Set("node_ids", []string{nodeID})

	return diags
}

// -------------------------------------------------------------------
// Read
// -------------------------------------------------------------------
func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", hostBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	} else if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to read host. HTTP status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "address")
	setStringField(d, result, "comment")
	setIntField(d, result, "platform")
	if zone, ok := result["zone"].(string); ok {
		d.Set("zone_id", zone)
	}
	if nodes, ok := result["nodes"].([]interface{}); ok {
		d.Set("node_ids", nodes)
	}

	if accounts, ok := result["accounts"].([]interface{}); ok {
		d.Set("accounts", flattenAccounts(accounts))
	}
	if protocols, ok := result["protocols"].([]interface{}); ok {
		d.Set("protocols", flattenProtocols(protocols))
	}

	return diags
}

// -------------------------------------------------------------------
// Update
// -------------------------------------------------------------------
func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	zoneID := d.Get("zone_id").(string)
	nodeIDsRaw := d.Get("node_ids").([]interface{})
	var nodeID string
	if len(nodeIDsRaw) > 0 {
		nodeID = nodeIDsRaw[0].(string)
	}

	if d.HasChange("zone_name") {
		newZoneName := d.Get("zone_name").(string)
		foundID, err := findZoneIDByName(c, newZoneName)
		if err != nil {
			return diag.FromErr(err)
		}
		zoneID = foundID
		d.Set("zone_id", foundID)
	}

	if d.HasChange("node_name") {
		newNodeName := d.Get("node_name").(string)
		foundID, err := findNodeIDByName(c, newNodeName)
		if err != nil {
			return diag.FromErr(err)
		}
		nodeID = foundID
		d.Set("node_ids", []string{foundID})
	}

	hostData := map[string]interface{}{
		"name":     d.Get("name").(string),
		"address":  d.Get("address").(string),
		"platform": d.Get("platform").(int),
		"zone":     zoneID,
		"nodes":    []string{nodeID},
	}

	if v, ok := d.GetOk("comment"); ok {
		hostData["comment"] = v.(string)
	}

	if v, ok := d.GetOk("accounts"); ok {
		hostData["accounts"] = expandAccounts(v.([]interface{}))
	}
	if v, ok := d.GetOk("protocols"); ok {
		hostData["protocols"] = expandProtocols(v.([]interface{}))
	}

	resp, err := c.doRequest("PUT", hostBasePath+d.Id()+"/", hostData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to update host. HTTP status: %d", resp.StatusCode)
	}

	return resourceHostRead(ctx, d, m)
}

// -------------------------------------------------------------------
// Delete
// -------------------------------------------------------------------
func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", hostBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Failed to delete host. HTTP status: %d", resp.StatusCode)
	}

	d.SetId("")
	return diags
}

// -------------------------------------------------------------------
// Get zone_id / node_id from zone_name / node_name
// -------------------------------------------------------------------
func findZoneIDByName(c *Config, zoneName string) (string, error) {
	resp, err := c.doRequest("GET", "assets/zones/", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to list zones, status=%d", resp.StatusCode)
	}

	var zones []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&zones); err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zoneNameAPI, ok := zone["name"].(string); ok {
			if strings.EqualFold(zoneNameAPI, zoneName) {
				if id, idOk := zone["id"].(string); idOk {
					return id, nil
				}
				return "", fmt.Errorf("zone '%s' found but has no 'id'", zoneName)
			}
		}
	}
	return "", fmt.Errorf("zone '%s' not found in JumpServer", zoneName)
}

func findNodeIDByName(c *Config, nodeName string) (string, error) {
	resp, err := c.doRequest("GET", "assets/nodes/", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to list nodes, status=%d", resp.StatusCode)
	}

	var nodes []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return "", err
	}

	for _, node := range nodes {
		if nodeNameAPI, ok := node["name"].(string); ok {
			if strings.EqualFold(nodeNameAPI, nodeName) {
				if id, idOk := node["id"].(string); idOk {
					return id, nil
				}
				return "", fmt.Errorf("node '%s' found but has no 'id'", nodeName)
			}
		}
	}
	return "", fmt.Errorf("node '%s' not found in JumpServer", nodeName)
}

func expandAccounts(list []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, item := range list {
		m := item.(map[string]interface{})
		acc := map[string]interface{}{
			"on_invalid":  m["on_invalid"].(string),
			"is_active":   m["is_active"].(bool),
			"name":        m["name"].(string),
			"username":    m["username"].(string),
			"secret_type": m["secret_type"].(string),
			"secret":      m["secret"].(string),
		}
		result = append(result, acc)
	}
	return result
}

func expandProtocols(list []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, item := range list {
		m := item.(map[string]interface{})
		proto := map[string]interface{}{
			"name": m["name"].(string),
			"port": m["port"].(int),
		}
		result = append(result, proto)
	}
	return result
}

func flattenAccounts(accounts []interface{}) []interface{} {
	var result []interface{}
	for _, a := range accounts {
		m := a.(map[string]interface{})
		acc := map[string]interface{}{
			"on_invalid":  m["on_invalid"],
			"is_active":   m["is_active"],
			"name":        m["name"],
			"username":    m["username"],
			"secret_type": m["secret_type"],
		}
		result = append(result, acc)
	}
	return result
}

func flattenProtocols(protocols []interface{}) []interface{} {
	var result []interface{}
	for _, p := range protocols {
		m := p.(map[string]interface{})
		proto := map[string]interface{}{
			"name": m["name"],
			"port": m["port"],
		}
		result = append(result, proto)
	}
	return result
}
