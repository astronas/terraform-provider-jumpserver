package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const sessionSharingBasePath = "terminal/session-sharings/"

func resourceSessionSharing() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSessionSharingCreate,
		ReadContext:   resourceSessionSharingRead,
		UpdateContext: resourceSessionSharingUpdate,
		DeleteContext: resourceSessionSharingDelete,

		Schema: map[string]*schema.Schema{
			"session": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The session UUID to share.",
			},
			"expired_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Expiration time in minutes (0 = never).",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The origin URI.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the sharing is active.",
			},
			"verify_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The verification code.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The sharing URL.",
			},
			"users_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display names of users.",
			},
			"action_permission": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action permission level.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created by user.",
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

func resourceSessionSharingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"session": d.Get("session").(string),
	}
	if v, ok := d.GetOk("expired_time"); ok {
		data["expired_time"] = v.(int)
	}
	if v, ok := d.GetOk("origin"); ok {
		data["origin"] = v.(string)
	}
	data["is_active"] = d.Get("is_active").(bool)

	resp, err := c.doRequest("POST", sessionSharingBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating session sharing: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving session sharing ID after creation, response: %v", result)
	}

	return resourceSessionSharingRead(ctx, d, m)
}

func resourceSessionSharingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", sessionSharingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading session sharing: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "session")
	setIntField(d, result, "expired_time")
	setStringField(d, result, "origin")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "verify_code")
	setStringField(d, result, "url")
	setStringField(d, result, "users_display")
	setStringField(d, result, "action_permission")
	setStringField(d, result, "created_by")
	setStringField(d, result, "org_id")
	setStringField(d, result, "org_name")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceSessionSharingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"session": d.Get("session").(string),
	}
	if v, ok := d.GetOk("expired_time"); ok {
		data["expired_time"] = v.(int)
	}
	if v, ok := d.GetOk("origin"); ok {
		data["origin"] = v.(string)
	}
	data["is_active"] = d.Get("is_active").(bool)

	resp, err := c.doRequest("PUT", sessionSharingBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating session sharing: %s", resp.Status)
	}

	return resourceSessionSharingRead(ctx, d, m)
}

func resourceSessionSharingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", sessionSharingBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting session sharing: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
