package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const commandGroupBasePath = "acls/command-groups/"

func resourceCommandGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommandGroupCreate,
		ReadContext:   resourceCommandGroupRead,
		UpdateContext: resourceCommandGroupUpdate,
		DeleteContext: resourceCommandGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the command group.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "command",
				ValidateFunc: validation.StringInSlice([]string{
					"command", "regex",
				}, false),
				Description: "Type of matching: command or regex.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Commands to match, one per line.",
			},
			"ignore_case": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to ignore case when matching commands.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the command group.",
			},
		},
	}
}

func resourceCommandGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"type":        d.Get("type").(string),
		"content":     d.Get("content").(string),
		"ignore_case": d.Get("ignore_case").(bool),
	}

	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", commandGroupBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_command_group", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating command group: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving command group ID after creation, response: %v", result)
	}

	return resourceCommandGroupRead(ctx, d, m)
}

func resourceCommandGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", commandGroupBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading command group: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "content")
	setBoolField(d, result, "ignore_case")
	setStringField(d, result, "comment")
	setEnumField(d, result, "type")

	return diags
}

func resourceCommandGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"type":        d.Get("type").(string),
		"content":     d.Get("content").(string),
		"ignore_case": d.Get("ignore_case").(bool),
		"comment":     d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", commandGroupBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating command group: %s", resp.Status)
	}

	return resourceCommandGroupRead(ctx, d, m)
}

func resourceCommandGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", commandGroupBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting command group: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
