package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const commandStorageBasePath = "terminal/command-storages/"

func resourceCommandStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommandStorageCreate,
		ReadContext:   resourceCommandStorageRead,
		UpdateContext: resourceCommandStorageUpdate,
		DeleteContext: resourceCommandStorageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the command storage.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"null", "server", "es",
				}, false),
				Description: "Storage type: null, server, es (Elasticsearch).",
			},
			"meta": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "{}",
				Description: "Storage configuration as JSON (e.g. Elasticsearch connection settings).",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is the default command storage.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the storage.",
			},
		},
	}
}

func resourceCommandStorageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var meta interface{}
	if err := json.Unmarshal([]byte(d.Get("meta").(string)), &meta); err != nil {
		return diag.Errorf("meta must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":       d.Get("name").(string),
		"type":       d.Get("type").(string),
		"meta":       meta,
		"is_default": d.Get("is_default").(bool),
		"comment":    d.Get("comment").(string),
	}

	resp, err := c.doRequest("POST", commandStorageBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_command_storage", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating command storage: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving command storage ID after creation, response: %v", result)
	}

	return resourceCommandStorageRead(ctx, d, m)
}

func resourceCommandStorageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", commandStorageBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading command storage: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setBoolField(d, result, "is_default")
	setStringField(d, result, "comment")
	setEnumField(d, result, "type")

	// meta: serialize back to JSON string
	if v, ok := result["meta"]; ok {
		b, _ := json.Marshal(v)
		d.Set("meta", string(b))
	}

	return diags
}

func resourceCommandStorageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var meta interface{}
	if err := json.Unmarshal([]byte(d.Get("meta").(string)), &meta); err != nil {
		return diag.Errorf("meta must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":       d.Get("name").(string),
		"type":       d.Get("type").(string),
		"meta":       meta,
		"is_default": d.Get("is_default").(bool),
		"comment":    d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", commandStorageBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating command storage: %s", resp.Status)
	}

	return resourceCommandStorageRead(ctx, d, m)
}

func resourceCommandStorageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", commandStorageBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting command storage: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
