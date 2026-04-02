package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const virtualAppBasePath = "terminal/virtual-apps/"

func resourceVirtualApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualAppCreate,
		ReadContext:   resourceVirtualAppRead,
		UpdateContext: resourceVirtualAppUpdate,
		DeleteContext: resourceVirtualAppDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the virtual app.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the virtual app.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Docker image name.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the virtual app.",
			},
			"author": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The author of the virtual app.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the virtual app is active.",
			},
			"image_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vnc",
				Description: "The protocol used by the image (vnc, rdp).",
			},
			"image_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5900,
				Description: "The port used by the image.",
			},
			"protocols": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of supported protocols.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of tags.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the virtual app.",
			},
		},
	}
}

func resourceVirtualAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"display_name":   d.Get("display_name").(string),
		"image_name":     d.Get("image_name").(string),
		"version":        d.Get("version").(string),
		"author":         d.Get("author").(string),
		"is_active":      d.Get("is_active").(bool),
		"image_protocol": d.Get("image_protocol").(string),
		"image_port":     d.Get("image_port").(int),
	}

	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}
	if v, ok := d.GetOk("protocols"); ok {
		var protocols interface{}
		if err := json.Unmarshal([]byte(v.(string)), &protocols); err != nil {
			return diag.Errorf("invalid protocols JSON: %s", err)
		}
		data["protocols"] = protocols
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags interface{}
		if err := json.Unmarshal([]byte(v.(string)), &tags); err != nil {
			return diag.Errorf("invalid tags JSON: %s", err)
		}
		data["tags"] = tags
	}

	resp, err := c.doRequest("POST", virtualAppBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_virtual_app", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating virtual app: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving virtual app ID after creation, response: %v", result)
	}

	return resourceVirtualAppRead(ctx, d, m)
}

func resourceVirtualAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", virtualAppBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading virtual app: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "display_name")
	setStringField(d, result, "image_name")
	setStringField(d, result, "version")
	setStringField(d, result, "author")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "image_protocol")
	setIntField(d, result, "image_port")
	setStringField(d, result, "comment")

	if v, ok := result["protocols"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("protocols", string(jsonBytes))
	}
	if v, ok := result["tags"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("tags", string(jsonBytes))
	}

	return diags
}

func resourceVirtualAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"display_name":   d.Get("display_name").(string),
		"image_name":     d.Get("image_name").(string),
		"version":        d.Get("version").(string),
		"author":         d.Get("author").(string),
		"is_active":      d.Get("is_active").(bool),
		"image_protocol": d.Get("image_protocol").(string),
		"image_port":     d.Get("image_port").(int),
		"comment":        d.Get("comment").(string),
	}

	if v, ok := d.GetOk("protocols"); ok {
		var protocols interface{}
		if err := json.Unmarshal([]byte(v.(string)), &protocols); err != nil {
			return diag.Errorf("invalid protocols JSON: %s", err)
		}
		data["protocols"] = protocols
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags interface{}
		if err := json.Unmarshal([]byte(v.(string)), &tags); err != nil {
			return diag.Errorf("invalid tags JSON: %s", err)
		}
		data["tags"] = tags
	}

	resp, err := c.doRequest("PUT", virtualAppBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating virtual app: %s", resp.Status)
	}

	return resourceVirtualAppRead(ctx, d, m)
}

func resourceVirtualAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", virtualAppBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting virtual app: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
