package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const nodeBasePath = "assets/nodes/"

func resourceNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeCreate,
		ReadContext:   resourceNodeRead,
		UpdateContext: resourceNodeUpdate,
		DeleteContext: resourceNodeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name (value) of the node.",
			},
			"full_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path of the node in the tree (e.g. /Default/Child).",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The tree key of the node (e.g. 1:1:2).",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the parent node. If omitted, the node is created under the default root.",
			},
			"child_mark": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The child mark counter for the node.",
			},
		},
	}
}

func createNodeRequest(c *Config, path string, nodeData map[string]interface{}) (string, diag.Diagnostics) {
	resp, err := c.doRequest("POST", path, nodeData)
	if err != nil {
		return "", diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", diag.Errorf("Error creating node: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", diag.FromErr(err)
	}

	id, ok := result["id"].(string)
	if !ok {
		return "", diag.Errorf("Error retrieving node ID after creation, response: %v", result)
	}
	return id, nil
}

func resourceNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	nodeData := map[string]interface{}{
		"value": d.Get("value").(string),
	}

	path := nodeBasePath
	if v, ok := d.GetOk("parent_id"); ok {
		path = nodeBasePath + v.(string) + "/children/"
	}

	id, diags := createNodeRequest(c, path, nodeData)
	if diags != nil {
		return diags
	}
	d.SetId(id)

	return resourceNodeRead(ctx, d, m)
}

func resourceNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", nodeBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading node: %s", resp.Status)
	}

	var node map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, node, "value")
	setStringField(d, node, "full_value")
	setStringField(d, node, "key")
	setIntField(d, node, "child_mark")

	return diags
}

func resourceNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	nodeData := map[string]interface{}{
		"value": d.Get("value").(string),
	}

	resp, err := c.doRequest("PUT", nodeBasePath+d.Id()+"/", nodeData)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating node: %s", resp.Status)
	}

	return resourceNodeRead(ctx, d, m)
}

func resourceNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", nodeBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting node: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
