package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const applyLoginTicketBasePath = "tickets/apply-login-tickets/"

func resourceApplyLoginTicket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplyLoginTicketCreate,
		ReadContext:   resourceApplyLoginTicketRead,
		UpdateContext: resourceApplyLoginTicketUpdate,
		DeleteContext: resourceApplyLoginTicketDelete,

		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the ticket.",
			},
			"apply_login_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The login IP address.",
			},
			"apply_login_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The login city.",
			},
			"apply_login_datetime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The login date and time.",
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

func resourceApplyLoginTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("apply_login_ip"); ok {
		data["apply_login_ip"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_city"); ok {
		data["apply_login_city"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_datetime"); ok {
		data["apply_login_datetime"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", applyLoginTicketBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating apply login ticket: %s", resp.Status)
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

	return resourceApplyLoginTicketRead(ctx, d, m)
}

func resourceApplyLoginTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", applyLoginTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading apply login ticket: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "title")
	setStringField(d, result, "apply_login_ip")
	setStringField(d, result, "apply_login_city")
	setStringField(d, result, "apply_login_datetime")
	setStringField(d, result, "comment")
	setEnumField(d, result, "state")
	setEnumField(d, result, "status")
	setStringField(d, result, "serial_num")
	setStringField(d, result, "org_name")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceApplyLoginTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("apply_login_ip"); ok {
		data["apply_login_ip"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_city"); ok {
		data["apply_login_city"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_datetime"); ok {
		data["apply_login_datetime"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", applyLoginTicketBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating apply login ticket: %s", resp.Status)
	}

	return resourceApplyLoginTicketRead(ctx, d, m)
}

func resourceApplyLoginTicketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", applyLoginTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting apply login ticket: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
