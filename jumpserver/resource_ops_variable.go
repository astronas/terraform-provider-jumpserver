package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const opsVariableBasePath = "ops/variables/"

func resourceOpsVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpsVariableCreate,
		ReadContext:   resourceOpsVariableRead,
		UpdateContext: resourceOpsVariableUpdate,
		DeleteContext: resourceOpsVariableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the variable.",
			},
			"var_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The variable name used in scripts.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "text",
				Description: "The type of the variable (text).",
			},
			"default_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The default value of the variable.",
			},
			"tips": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Help text for the variable.",
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the variable is required.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the variable.",
			},
		},
	}
}

func resourceOpsVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":          d.Get("name").(string),
		"var_name":      d.Get("var_name").(string),
		"type":          d.Get("type").(string),
		"default_value": d.Get("default_value").(string),
		"tips":          d.Get("tips").(string),
		"required":      d.Get("required").(bool),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", opsVariableBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_ops_variable", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating ops variable: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving ops variable ID after creation, response: %v", result)
	}

	return resourceOpsVariableRead(ctx, d, m)
}

func resourceOpsVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", opsVariableBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading ops variable: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "var_name")
	setStringField(d, result, "type")
	setStringField(d, result, "default_value")
	setStringField(d, result, "tips")
	setBoolField(d, result, "required")
	setStringField(d, result, "comment")

	return diags
}

func resourceOpsVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":          d.Get("name").(string),
		"var_name":      d.Get("var_name").(string),
		"type":          d.Get("type").(string),
		"default_value": d.Get("default_value").(string),
		"tips":          d.Get("tips").(string),
		"required":      d.Get("required").(bool),
		"comment":       d.Get("comment").(string),
	}

	resp, err := c.doRequest("PUT", opsVariableBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating ops variable: %s", resp.Status)
	}

	return resourceOpsVariableRead(ctx, d, m)
}

func resourceOpsVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", opsVariableBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting ops variable: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
