package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ticketBasePath = "tickets/tickets/"

func resourceTicket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTicketCreate,
		ReadContext:   resourceTicketRead,
		UpdateContext: resourceTicketUpdate,
		DeleteContext: resourceTicketDelete,

		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the ticket.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A comment or description.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the ticket.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the ticket (pending, approved, rejected, closed).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the ticket (open, closed).",
			},
			"serial_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The serial number of the ticket.",
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

func resourceTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", ticketBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating ticket: %s", resp.Status)
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

	return resourceTicketRead(ctx, d, m)
}

func resourceTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", ticketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading ticket: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "title")
	setStringField(d, result, "comment")
	setEnumField(d, result, "type")
	setEnumField(d, result, "state")
	setEnumField(d, result, "status")
	setStringField(d, result, "serial_num")
	setStringField(d, result, "org_id")
	setStringField(d, result, "org_name")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", ticketBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating ticket: %s", resp.Status)
	}

	return resourceTicketRead(ctx, d, m)
}

func resourceTicketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", ticketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting ticket: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
