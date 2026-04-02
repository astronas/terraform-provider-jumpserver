package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const systemMsgSubscriptionBasePath = "notifications/system-msg-subscription/"

func resourceSystemMsgSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemMsgSubscriptionCreate,
		ReadContext:   resourceSystemMsgSubscriptionRead,
		UpdateContext: resourceSystemMsgSubscriptionUpdate,
		DeleteContext: resourceSystemMsgSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"message_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The message type identifier (used as the resource ID).",
			},
			"message_type_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable label for the message type.",
			},
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user UUIDs to subscribe.",
			},
			"groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user group UUIDs to subscribe.",
			},
			"receive_backends": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of notification backend names.",
			},
			"receivers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The resolved list of receiver display names.",
			},
		},
	}
}

func systemMsgSubscriptionBuildPayload(d *schema.ResourceData) map[string]interface{} {
	data := map[string]interface{}{}

	if v, ok := d.GetOk("users"); ok {
		data["users"] = v.([]interface{})
	} else {
		data["users"] = []interface{}{}
	}
	if v, ok := d.GetOk("groups"); ok {
		data["groups"] = v.([]interface{})
	} else {
		data["groups"] = []interface{}{}
	}
	if v, ok := d.GetOk("receive_backends"); ok {
		data["receive_backends"] = v.([]interface{})
	} else {
		data["receive_backends"] = []interface{}{}
	}

	return data
}

func resourceSystemMsgSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	messageType := d.Get("message_type").(string)
	d.SetId(messageType)

	return resourceSystemMsgSubscriptionUpdate(ctx, d, m)
}

func resourceSystemMsgSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", systemMsgSubscriptionBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading system message subscription: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "message_type")
	setStringField(d, result, "message_type_label")
	setObjectIDsField(d, result, "users")
	setObjectIDsField(d, result, "groups")
	if v, ok := result["receive_backends"].([]interface{}); ok {
		d.Set("receive_backends", v)
	}
	if v, ok := result["receivers"].([]interface{}); ok {
		d.Set("receivers", v)
	}

	return diags
}

func resourceSystemMsgSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := systemMsgSubscriptionBuildPayload(d)

	resp, err := c.doRequest("PUT", systemMsgSubscriptionBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating system message subscription: %s", resp.Status)
	}

	return resourceSystemMsgSubscriptionRead(ctx, d, m)
}

func resourceSystemMsgSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"users":            []interface{}{},
		"groups":           []interface{}{},
		"receive_backends": []interface{}{},
	}

	resp, err := c.doRequest("PUT", systemMsgSubscriptionBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	d.SetId("")
	return nil
}
