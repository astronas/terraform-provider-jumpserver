package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const ticketFlowBasePath = "tickets/flows/"

func resourceTicketFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTicketFlowCreate,
		ReadContext:   resourceTicketFlowRead,
		UpdateContext: resourceTicketFlowUpdate,
		DeleteContext: resourceTicketFlowDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"general", "apply_asset", "login_confirm",
					"command_confirm", "login_asset_confirm",
				}, false),
				Description: "The type of the ticket flow.",
			},
			"approval_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 3),
				Description:  "The number of approval levels (1-3).",
			},
			"rules": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON string of approval rules.",
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
		},
	}
}

func ticketFlowBuildPayload(d *schema.ResourceData) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"type":           map[string]interface{}{"value": d.Get("type").(string)},
		"approval_level": d.Get("approval_level").(int),
	}

	var rules interface{}
	if err := json.Unmarshal([]byte(d.Get("rules").(string)), &rules); err != nil {
		return nil, err
	}
	data["rules"] = rules

	return data, nil
}

func resourceTicketFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data, err := ticketFlowBuildPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.doRequest("POST", ticketFlowBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating ticket flow: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving ticket flow ID after creation, response: %v", result)
	}

	return resourceTicketFlowRead(ctx, d, m)
}

func resourceTicketFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", ticketFlowBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading ticket flow: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setEnumField(d, result, "type")
	setIntField(d, result, "approval_level")
	setStringField(d, result, "org_id")
	setStringField(d, result, "org_name")

	if v, ok := result["rules"]; ok {
		jsonBytes, _ := json.Marshal(v)
		d.Set("rules", string(jsonBytes))
	}

	return diags
}

func resourceTicketFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data, err := ticketFlowBuildPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.doRequest("PUT", ticketFlowBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating ticket flow: %s", resp.Status)
	}

	return resourceTicketFlowRead(ctx, d, m)
}

func resourceTicketFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", ticketFlowBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting ticket flow: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
