package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const applyAssetTicketBasePath = "tickets/apply-asset-tickets/"

func resourceApplyAssetTicket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplyAssetTicketCreate,
		ReadContext:   resourceApplyAssetTicketRead,
		UpdateContext: resourceApplyAssetTicketUpdate,
		DeleteContext: resourceApplyAssetTicketDelete,

		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the ticket.",
			},
			"apply_nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of node UUIDs to apply for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"apply_assets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of asset UUIDs to apply for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"apply_accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of accounts to apply for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"apply_actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of actions to apply for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"apply_date_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start date for the applied permission.",
			},
			"apply_date_expired": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration date for the applied permission.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A comment or description.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the ticket.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the ticket.",
			},
			"serial_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The serial number of the ticket.",
			},
			"apply_permission_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the applied permission.",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization name.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date.",
			},
		},
	}
}

func resourceApplyAssetTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("apply_nodes"); ok {
		data["apply_nodes"] = v.([]interface{})
	}
	if v, ok := d.GetOk("apply_assets"); ok {
		data["apply_assets"] = v.([]interface{})
	}
	if v, ok := d.GetOk("apply_accounts"); ok {
		data["apply_accounts"] = v.([]interface{})
	}
	if v, ok := d.GetOk("apply_actions"); ok {
		data["apply_actions"] = v.([]interface{})
	}
	if v, ok := d.GetOk("apply_date_start"); ok {
		data["apply_date_start"] = v.(string)
	}
	if v, ok := d.GetOk("apply_date_expired"); ok {
		data["apply_date_expired"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", applyAssetTicketBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating apply asset ticket: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving ticket ID after creation, response: %v", result)
	}

	return resourceApplyAssetTicketRead(ctx, d, m)
}

func resourceApplyAssetTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", applyAssetTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading apply asset ticket: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "title")
	setStringField(d, result, "comment")
	setEnumField(d, result, "state")
	setEnumField(d, result, "status")
	setStringField(d, result, "serial_num")
	setStringField(d, result, "apply_permission_name")
	setStringField(d, result, "org_name")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceApplyAssetTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", applyAssetTicketBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating apply asset ticket: %s", resp.Status)
	}

	return resourceApplyAssetTicketRead(ctx, d, m)
}

func resourceApplyAssetTicketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", applyAssetTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting apply asset ticket: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
