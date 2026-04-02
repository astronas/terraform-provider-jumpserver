package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const applyCommandTicketBasePath = "tickets/apply-command-tickets/"

func resourceApplyCommandTicket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplyCommandTicketCreate,
		ReadContext:   resourceApplyCommandTicketRead,
		UpdateContext: resourceApplyCommandTicketUpdate,
		DeleteContext: resourceApplyCommandTicketDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the ticket.",
			},
			"apply_run_asset": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The asset to run the command on.",
			},
			"apply_run_account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account to run the command as.",
			},
			"apply_run_command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The command to run.",
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
				Description: "The serial number.",
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

func resourceApplyCommandTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title":             d.Get("title").(string),
		"apply_run_asset":   d.Get("apply_run_asset").(string),
		"apply_run_account": d.Get("apply_run_account").(string),
		"apply_run_command": d.Get("apply_run_command").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", applyCommandTicketBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_apply_command_ticket", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating apply command ticket: %s", resp.Status)
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

	return resourceApplyCommandTicketRead(ctx, d, m)
}

func resourceApplyCommandTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", applyCommandTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading apply command ticket: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "title")
	setStringField(d, result, "apply_run_asset")
	setStringField(d, result, "apply_run_account")
	setStringField(d, result, "apply_run_command")
	setStringField(d, result, "comment")
	setEnumField(d, result, "state")
	setEnumField(d, result, "status")
	setStringField(d, result, "serial_num")
	setStringField(d, result, "org_name")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceApplyCommandTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title":             d.Get("title").(string),
		"apply_run_asset":   d.Get("apply_run_asset").(string),
		"apply_run_account": d.Get("apply_run_account").(string),
		"apply_run_command": d.Get("apply_run_command").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", applyCommandTicketBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating apply command ticket: %s", resp.Status)
	}

	return resourceApplyCommandTicketRead(ctx, d, m)
}

func resourceApplyCommandTicketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", applyCommandTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting apply command ticket: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
