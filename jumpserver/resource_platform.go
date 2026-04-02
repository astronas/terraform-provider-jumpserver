package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const platformBasePath = "assets/platforms/"

func resourcePlatform() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlatformCreate,
		ReadContext:   resourcePlatformRead,
		UpdateContext: resourcePlatformUpdate,
		DeleteContext: resourcePlatformDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the platform.",
			},
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The category of the platform (host, device, database, cloud, web, gpt, custom).",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the platform.",
			},
			"charset": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "utf-8",
				Description: "The character encoding (utf-8 or gbk).",
			},
			"domain_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether domain (zone) is enabled for this platform.",
			},
			"su_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether su (switch user) is enabled.",
			},
			"su_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The su method (su, sudo, etc.).",
			},
			"protocols": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of protocols configuration.",
			},
			"automation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of automation configuration.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the platform.",
			},
		},
	}
}

func platformBuildPayload(d *schema.ResourceData) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"category":       map[string]interface{}{"value": d.Get("category").(string)},
		"charset":        map[string]interface{}{"value": d.Get("charset").(string)},
		"domain_enabled": d.Get("domain_enabled").(bool),
		"su_enabled":     d.Get("su_enabled").(bool),
	}

	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}
	if v, ok := d.GetOk("su_method"); ok {
		data["su_method"] = map[string]interface{}{"value": v.(string)}
	}

	if v, ok := d.GetOk("protocols"); ok {
		var protocols interface{}
		if err := json.Unmarshal([]byte(v.(string)), &protocols); err != nil {
			return nil, fmt.Errorf("invalid protocols JSON: %s", err)
		}
		data["protocols"] = protocols
	}
	if v, ok := d.GetOk("automation"); ok {
		var automation interface{}
		if err := json.Unmarshal([]byte(v.(string)), &automation); err != nil {
			return nil, fmt.Errorf("invalid automation JSON: %s", err)
		}
		data["automation"] = automation
	}

	return data, nil
}

func resourcePlatformCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data, err := platformBuildPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.doRequest("POST", platformBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_platform", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error creating platform: %s", resp.Status)
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
		return diag.Errorf("Error retrieving platform ID after creation, response: %v", result)
	}

	return resourcePlatformRead(ctx, d, m)
}

func resourcePlatformRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", platformBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading platform: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "comment")
	setEnumField(d, result, "category")
	setEnumField(d, result, "type")
	setEnumField(d, result, "charset")
	setBoolField(d, result, "domain_enabled")
	setBoolField(d, result, "su_enabled")
	setEnumField(d, result, "su_method")

	if v, ok := result["protocols"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("protocols", string(jsonBytes))
	}
	if v, ok := result["automation"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("automation", string(jsonBytes))
	}

	return diags
}

func resourcePlatformUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data, err := platformBuildPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.doRequest("PUT", platformBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating platform: %s", resp.Status)
	}

	return resourcePlatformRead(ctx, d, m)
}

func resourcePlatformDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", platformBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting platform: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
