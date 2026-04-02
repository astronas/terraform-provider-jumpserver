package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const playbookBasePath = "ops/playbooks/"

func resourcePlaybook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookCreate,
		ReadContext:   resourcePlaybookRead,
		UpdateContext: resourcePlaybookUpdate,
		DeleteContext: resourcePlaybookDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the playbook.",
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public", "private",
				}, false),
				Description: "The scope of the playbook: public or private.",
			},
			"create_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "blank",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"blank", "vcs",
				}, false),
				Description: "The creation method: blank or vcs.",
			},
			"vcs_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VCS repository URL (used when create_method is vcs).",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the playbook.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path to the playbook on the server.",
			},
		},
	}
}

func resourcePlaybookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":          d.Get("name").(string),
		"scope":         d.Get("scope").(string),
		"create_method": d.Get("create_method").(string),
		"comment":       d.Get("comment").(string),
	}

	if v, ok := d.GetOk("vcs_url"); ok {
		data["vcs_url"] = v.(string)
	}

	resp, err := c.doRequest("POST", playbookBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_playbook", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating playbook: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving playbook ID after creation, response: %v", result)
	}

	return resourcePlaybookRead(ctx, d, m)
}

func resourcePlaybookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", playbookBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading playbook: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "path")
	setStringField(d, result, "comment")
	setStringField(d, result, "vcs_url")
	setEnumField(d, result, "scope")
	setEnumField(d, result, "create_method")

	return diags
}

func resourcePlaybookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":    d.Get("name").(string),
		"scope":   d.Get("scope").(string),
		"comment": d.Get("comment").(string),
	}

	if v, ok := d.GetOk("vcs_url"); ok {
		data["vcs_url"] = v.(string)
	}

	resp, err := c.doRequest("PUT", playbookBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating playbook: %s", resp.Status)
	}

	return resourcePlaybookRead(ctx, d, m)
}

func resourcePlaybookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", playbookBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting playbook: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
