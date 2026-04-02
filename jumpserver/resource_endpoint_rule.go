package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const endpointRuleBasePath = "terminal/endpoint-rules/"

func resourceEndpointRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointRuleCreate,
		ReadContext:   resourceEndpointRuleRead,
		UpdateContext: resourceEndpointRuleUpdate,
		DeleteContext: resourceEndpointRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the endpoint rule.",
			},
			"ip_group": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IP ranges to match (e.g. [\"192.168.1.0/24\", \"10.0.0.1-10.0.0.50\", \"*\"]). Defaults to [\"*\"].",
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validation.IntBetween(1, 100),
				Description:  "Priority (1-100). Lower values match first.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Endpoint UUID to route matching traffic to.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the rule is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the rule.",
			},
		},
	}
}

func resourceEndpointRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"priority":  d.Get("priority").(int),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("ip_group"); ok {
		data["ip_group"] = v.([]interface{})
	} else {
		data["ip_group"] = []string{"*"}
	}

	if v, ok := d.GetOk("endpoint"); ok {
		data["endpoint"] = v.(string)
	}

	resp, err := c.doRequest("POST", endpointRuleBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_endpoint_rule", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating endpoint rule: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving endpoint rule ID after creation, response: %v", result)
	}

	return resourceEndpointRuleRead(ctx, d, m)
}

func resourceEndpointRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", endpointRuleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading endpoint rule: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setIntField(d, result, "priority")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")

	if v, ok := result["ip_group"].([]interface{}); ok {
		d.Set("ip_group", v)
	}

	setObjectIDField(d, result, "endpoint")

	return diags
}

func resourceEndpointRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"priority":  d.Get("priority").(int),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}

	if v, ok := d.GetOk("ip_group"); ok {
		data["ip_group"] = v.([]interface{})
	} else {
		data["ip_group"] = []string{"*"}
	}

	if v, ok := d.GetOk("endpoint"); ok {
		data["endpoint"] = v.(string)
	}

	resp, err := c.doRequest("PUT", endpointRuleBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating endpoint rule: %s", resp.Status)
	}

	return resourceEndpointRuleRead(ctx, d, m)
}

func resourceEndpointRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", endpointRuleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting endpoint rule: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
