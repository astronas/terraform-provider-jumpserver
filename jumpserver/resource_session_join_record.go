package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const sessionJoinRecordBasePath = "terminal/session-join-records/"

func resourceSessionJoinRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSessionJoinRecordCreate,
		ReadContext:   resourceSessionJoinRecordRead,
		UpdateContext: resourceSessionJoinRecordUpdate,
		DeleteContext: resourceSessionJoinRecordDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"session": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The session UUID.",
			},
			"sharing": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The session sharing UUID.",
			},
			"joiner": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The joiner user UUID.",
			},
			"verify_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The verification code.",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The remote address.",
			},
			"login_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "WT",
				Description: "The login source: ST, RT, WT.",
			},
			"reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The reason for joining.",
			},
			"is_success": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the join was successful.",
			},
			"is_finished": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the join session is finished.",
			},
			"joiner_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The joiner display name.",
			},
			"action_permission": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action permission.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created by user.",
			},
			"date_joined": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The join date.",
			},
			"date_left": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The leave date.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
		},
	}
}

func resourceSessionJoinRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"session":     d.Get("session").(string),
		"sharing":     d.Get("sharing").(string),
		"verify_code": d.Get("verify_code").(string),
	}
	if v, ok := d.GetOk("joiner"); ok {
		data["joiner"] = v.(string)
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	if v, ok := d.GetOk("login_from"); ok {
		data["login_from"] = v.(string)
	}
	if v, ok := d.GetOk("reason"); ok {
		data["reason"] = v.(string)
	}
	data["is_success"] = d.Get("is_success").(bool)
	data["is_finished"] = d.Get("is_finished").(bool)

	resp, err := c.doRequest("POST", sessionJoinRecordBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating session join record: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving session join record ID after creation, response: %v", result)
	}

	return resourceSessionJoinRecordRead(ctx, d, m)
}

func resourceSessionJoinRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", sessionJoinRecordBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading session join record: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "session")
	setObjectIDField(d, result, "sharing")
	setObjectIDField(d, result, "joiner")
	setStringField(d, result, "verify_code")
	setStringField(d, result, "remote_addr")
	setEnumField(d, result, "login_from")
	setStringField(d, result, "reason")
	setBoolField(d, result, "is_success")
	setBoolField(d, result, "is_finished")
	setStringField(d, result, "joiner_display")
	setStringField(d, result, "action_permission")
	setStringField(d, result, "created_by")
	setStringField(d, result, "date_joined")
	setStringField(d, result, "date_left")
	setStringField(d, result, "date_created")

	return diags
}

func resourceSessionJoinRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"session":     d.Get("session").(string),
		"sharing":     d.Get("sharing").(string),
		"verify_code": d.Get("verify_code").(string),
	}
	if v, ok := d.GetOk("joiner"); ok {
		data["joiner"] = v.(string)
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	if v, ok := d.GetOk("login_from"); ok {
		data["login_from"] = v.(string)
	}
	if v, ok := d.GetOk("reason"); ok {
		data["reason"] = v.(string)
	}
	data["is_success"] = d.Get("is_success").(bool)
	data["is_finished"] = d.Get("is_finished").(bool)

	resp, err := c.doRequest("PUT", sessionJoinRecordBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating session join record: %s", resp.Status)
	}

	return resourceSessionJoinRecordRead(ctx, d, m)
}

func resourceSessionJoinRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", sessionJoinRecordBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting session join record: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
