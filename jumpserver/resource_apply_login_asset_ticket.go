package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const applyLoginAssetTicketBasePath = "tickets/apply-login-asset-tickets/"

func resourceApplyLoginAssetTicket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplyLoginAssetTicketCreate,
		ReadContext:   resourceApplyLoginAssetTicketRead,
		UpdateContext: resourceApplyLoginAssetTicketUpdate,
		DeleteContext: resourceApplyLoginAssetTicketDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the ticket.",
			},
			"apply_login_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the login user.",
			},
			"apply_login_asset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the login asset.",
			},
			"apply_login_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The login account name.",
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

func resourceApplyLoginAssetTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("apply_login_user"); ok {
		data["apply_login_user"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_asset"); ok {
		data["apply_login_asset"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_account"); ok {
		data["apply_login_account"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", applyLoginAssetTicketBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_apply_login_asset_ticket", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating apply login asset ticket: %s", resp.Status)
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

	return resourceApplyLoginAssetTicketRead(ctx, d, m)
}

func resourceApplyLoginAssetTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", applyLoginAssetTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading apply login asset ticket: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "title")
	setObjectIDField(d, result, "apply_login_user")
	setObjectIDField(d, result, "apply_login_asset")
	setStringField(d, result, "apply_login_account")
	setStringField(d, result, "comment")
	setEnumField(d, result, "state")
	setEnumField(d, result, "status")
	setStringField(d, result, "serial_num")
	setStringField(d, result, "org_name")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceApplyLoginAssetTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if v, ok := d.GetOk("apply_login_user"); ok {
		data["apply_login_user"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_asset"); ok {
		data["apply_login_asset"] = v.(string)
	}
	if v, ok := d.GetOk("apply_login_account"); ok {
		data["apply_login_account"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", applyLoginAssetTicketBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating apply login asset ticket: %s", resp.Status)
	}

	return resourceApplyLoginAssetTicketRead(ctx, d, m)
}

func resourceApplyLoginAssetTicketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", applyLoginAssetTicketBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting apply login asset ticket: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
