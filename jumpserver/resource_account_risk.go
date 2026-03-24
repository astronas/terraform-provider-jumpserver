package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const accountRiskBasePath = "accounts/account-risks/"

func resourceAccountRisk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountRiskCreate,
		ReadContext:   resourceAccountRiskRead,
		UpdateContext: resourceAccountRiskUpdate,
		DeleteContext: resourceAccountRiskDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username associated with the risk entry.",
			},
			"asset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the asset associated with the risk.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "pending",
				Description: "The status of the risk (pending, confirmed, ignored).",
			},
			"details": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Additional details about the account risk.",
			},
			"risk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The computed risk level (JSON).",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
		},
	}
}

func resourceAccountRiskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"username": d.Get("username").(string),
	}
	if v, ok := d.GetOk("asset"); ok {
		data["asset"] = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		data["status"] = v.(string)
	}
	if v, ok := d.GetOk("details"); ok {
		data["details"] = v.(string)
	}

	resp, err := c.doRequest("POST", accountRiskBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating account risk: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving account risk ID after creation, response: %v", result)
	}

	return resourceAccountRiskRead(ctx, d, m)
}

func resourceAccountRiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", accountRiskBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading account risk: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "username")
	setObjectIDField(d, result, "asset")
	setStringField(d, result, "status")
	setStringField(d, result, "details")
	setStringField(d, result, "date_created")

	if risk, ok := result["risk"]; ok {
		riskJSON, _ := json.Marshal(risk)
		d.Set("risk", string(riskJSON))
	}

	return diags
}

func resourceAccountRiskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"username": d.Get("username").(string),
	}
	if v, ok := d.GetOk("asset"); ok {
		data["asset"] = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		data["status"] = v.(string)
	}
	if v, ok := d.GetOk("details"); ok {
		data["details"] = v.(string)
	}

	resp, err := c.doRequest("PUT", accountRiskBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating account risk: %s", resp.Status)
	}

	return resourceAccountRiskRead(ctx, d, m)
}

func resourceAccountRiskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", accountRiskBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting account risk: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
