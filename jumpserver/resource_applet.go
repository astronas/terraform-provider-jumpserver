package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const appletBasePath = "terminal/applets/"

func resourceApplet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppletCreate,
		ReadContext:   resourceAppletRead,
		UpdateContext: resourceAppletUpdate,
		DeleteContext: resourceAppletDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the applet.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the applet.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the applet.",
			},
			"author": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The author of the applet.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "general",
				Description: "The type of the applet (general, web).",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the applet is active.",
			},
			"can_concurrent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the applet supports concurrent sessions.",
			},
			"protocols": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of supported protocols.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string of applet tags.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the applet.",
			},
			"edition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The edition of the applet (community, enterprise).",
			},
		},
	}
}

func resourceAppletCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"display_name":   d.Get("display_name").(string),
		"version":        d.Get("version").(string),
		"author":         d.Get("author").(string),
		"type":           d.Get("type").(string),
		"is_active":      d.Get("is_active").(bool),
		"can_concurrent": d.Get("can_concurrent").(bool),
	}

	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}
	if v, ok := d.GetOk("edition"); ok {
		data["edition"] = v.(string)
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

	resp, err := c.doRequest("POST", appletBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_applet", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating applet: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving applet ID after creation, response: %v", result)
	}

	return resourceAppletRead(ctx, d, m)
}

func resourceAppletRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", appletBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading applet: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "display_name")
	setStringField(d, result, "version")
	setStringField(d, result, "author")
	setEnumField(d, result, "type")
	setBoolField(d, result, "is_active")
	setBoolField(d, result, "can_concurrent")
	setStringField(d, result, "comment")
	setStringField(d, result, "edition")

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

func resourceAppletUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"display_name":   d.Get("display_name").(string),
		"version":        d.Get("version").(string),
		"author":         d.Get("author").(string),
		"type":           d.Get("type").(string),
		"is_active":      d.Get("is_active").(bool),
		"can_concurrent": d.Get("can_concurrent").(bool),
		"comment":        d.Get("comment").(string),
	}

	if v, ok := d.GetOk("edition"); ok {
		data["edition"] = v.(string)
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

	resp, err := c.doRequest("PUT", appletBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating applet: %s", resp.Status)
	}

	return resourceAppletRead(ctx, d, m)
}

func resourceAppletDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", appletBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting applet: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
