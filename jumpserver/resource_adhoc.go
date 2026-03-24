package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const adhocBasePath = "ops/adhocs/"

func resourceAdHoc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdHocCreate,
		ReadContext:   resourceAdHocRead,
		UpdateContext: resourceAdHocUpdate,
		DeleteContext: resourceAdHocDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ad-hoc command.",
			},
			"module": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shell",
				ValidateFunc: validation.StringInSlice([]string{
					"shell", "winshell", "python", "raw",
				}, false),
				Description: "The module to use for execution.",
			},
			"args": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The command arguments.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the ad-hoc command.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization ID.",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization name.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created this ad-hoc command.",
			},
		},
	}
}

func resourceAdHocCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":   d.Get("name").(string),
		"module": d.Get("module").(string),
		"args":   d.Get("args").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", adhocBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating ad-hoc command: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving ad-hoc command ID after creation, response: %v", result)
	}

	return resourceAdHocRead(ctx, d, m)
}

func resourceAdHocRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", adhocBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading ad-hoc command: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "module")
	setStringField(d, result, "args")
	setStringField(d, result, "comment")
	setStringField(d, result, "org_id")
	setStringField(d, result, "org_name")
	setStringField(d, result, "created_by")

	return diags
}

func resourceAdHocUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":    d.Get("name").(string),
		"module":  d.Get("module").(string),
		"args":    d.Get("args").(string),
		"comment": d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", adhocBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating ad-hoc command: %s", resp.Status)
	}

	return resourceAdHocRead(ctx, d, m)
}

func resourceAdHocDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", adhocBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting ad-hoc command: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
