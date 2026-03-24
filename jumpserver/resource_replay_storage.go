package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const replayStorageBasePath = "terminal/replay-storages/"

func resourceReplayStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReplayStorageCreate,
		ReadContext:   resourceReplayStorageRead,
		UpdateContext: resourceReplayStorageUpdate,
		DeleteContext: resourceReplayStorageDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the replay storage.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"null", "server", "s3", "ceph", "swift",
					"oss", "azure", "obs", "cos", "sftp",
				}, false),
				Description: "Storage type: null, server, s3, ceph, swift, oss, azure, obs, cos, sftp.",
			},
			"meta": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "{}",
				Description: "Storage configuration as JSON (e.g. S3 bucket, credentials, region).",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is the default replay storage.",
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

func resourceReplayStorageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resp, err := c.doRequest("POST", replayStorageBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating replay storage: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving replay storage ID after creation, response: %v", result)
	}

	return resourceReplayStorageRead(ctx, d, m)
}

func resourceReplayStorageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", replayStorageBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading replay storage: %s", resp.Status)
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

func resourceReplayStorageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resp, err := c.doRequest("PUT", replayStorageBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating replay storage: %s", resp.Status)
	}

	return resourceReplayStorageRead(ctx, d, m)
}

func resourceReplayStorageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", replayStorageBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting replay storage: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
