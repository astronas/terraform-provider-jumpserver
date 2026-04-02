package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const sessionBasePath = "terminal/sessions/"

func resourceSession() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSessionCreate,
		ReadContext:   resourceSessionRead,
		UpdateContext: resourceSessionUpdate,
		DeleteContext: resourceSessionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username for the session.",
			},
			"asset": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The asset name for the session.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user UUID.",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The asset UUID.",
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account name.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account ID.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol (e.g. ssh, rdp).",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "normal",
				Description: "The session type: normal, tunnel, command.",
			},
			"login_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ST",
				Description: "Login source: ST, RT, WT.",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The remote address.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A comment.",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session is locked.",
			},
			"is_finished": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session is finished.",
			},
			"is_success": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session login was successful.",
			},
			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The session duration in seconds.",
			},
			"command_amount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of commands executed.",
			},
			"error_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The error reason if login failed.",
			},
			"terminal_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The terminal display name.",
			},
			"can_replay": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session can be replayed.",
			},
			"can_join": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session can be joined.",
			},
			"can_terminate": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session can be terminated.",
			},
			"date_start": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session start date.",
			},
			"date_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session end date.",
			},
		},
	}
}

func resourceSessionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":       d.Get("user").(string),
		"asset":      d.Get("asset").(string),
		"user_id":    d.Get("user_id").(string),
		"asset_id":   d.Get("asset_id").(string),
		"account":    d.Get("account").(string),
		"account_id": d.Get("account_id").(string),
		"protocol":   d.Get("protocol").(string),
	}
	if v, ok := d.GetOk("type"); ok {
		data["type"] = v.(string)
	}
	if v, ok := d.GetOk("login_from"); ok {
		data["login_from"] = v.(string)
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", sessionBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating session: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving session ID after creation, response: %v", result)
	}

	return resourceSessionRead(ctx, d, m)
}

func resourceSessionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", sessionBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading session: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "user")
	setStringField(d, result, "asset")
	setStringField(d, result, "user_id")
	setStringField(d, result, "asset_id")
	setStringField(d, result, "account")
	setStringField(d, result, "account_id")
	setStringField(d, result, "protocol")
	setEnumField(d, result, "type")
	setEnumField(d, result, "login_from")
	setStringField(d, result, "remote_addr")
	setStringField(d, result, "comment")
	setBoolField(d, result, "is_locked")
	setBoolField(d, result, "is_finished")
	setBoolField(d, result, "is_success")
	setIntField(d, result, "duration")
	setIntField(d, result, "command_amount")
	setStringField(d, result, "error_reason")
	setStringField(d, result, "terminal_display")
	setBoolField(d, result, "can_replay")
	setBoolField(d, result, "can_join")
	setBoolField(d, result, "can_terminate")
	setStringField(d, result, "date_start")
	setStringField(d, result, "date_end")

	return diags
}

func resourceSessionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":       d.Get("user").(string),
		"asset":      d.Get("asset").(string),
		"user_id":    d.Get("user_id").(string),
		"asset_id":   d.Get("asset_id").(string),
		"account":    d.Get("account").(string),
		"account_id": d.Get("account_id").(string),
		"protocol":   d.Get("protocol").(string),
	}
	if v, ok := d.GetOk("type"); ok {
		data["type"] = v.(string)
	}
	if v, ok := d.GetOk("login_from"); ok {
		data["login_from"] = v.(string)
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", sessionBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating session: %s", resp.Status)
	}

	return resourceSessionRead(ctx, d, m)
}

func resourceSessionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", sessionBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting session: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
