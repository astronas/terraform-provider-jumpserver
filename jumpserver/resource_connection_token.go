package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const connectionTokenBasePath = "authentication/connection-token/"

func resourceConnectionToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionTokenCreate,
		ReadContext:   resourceConnectionTokenRead,
		UpdateContext: resourceConnectionTokenUpdate,
		DeleteContext: resourceConnectionTokenDelete,

		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account name.",
			},
			"connect_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The connection method.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ssh",
				Description: "The protocol (e.g. ssh, rdp).",
			},
			"asset": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The asset UUID.",
			},
			"input_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The input username.",
			},
			"input_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "The input secret (write-only).",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The remote address.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the token is active.",
			},
			"is_reusable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the token is reusable.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The token value.",
			},
			"user_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user display name.",
			},
			"asset_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset display name.",
			},
			"is_expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the token is expired.",
			},
			"expire_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration time in seconds.",
			},
			"date_expired": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date.",
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

func resourceConnectionTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"account":        d.Get("account").(string),
		"connect_method": d.Get("connect_method").(string),
	}
	if v, ok := d.GetOk("protocol"); ok {
		data["protocol"] = v.(string)
	}
	if v, ok := d.GetOk("asset"); ok {
		data["asset"] = v.(string)
	}
	if v, ok := d.GetOk("input_username"); ok {
		data["input_username"] = v.(string)
	}
	if v, ok := d.GetOk("input_secret"); ok {
		data["input_secret"] = v.(string)
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	data["is_active"] = d.Get("is_active").(bool)
	data["is_reusable"] = d.Get("is_reusable").(bool)

	resp, err := c.doRequest("POST", connectionTokenBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating connection token: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving connection token ID after creation, response: %v", result)
	}

	return resourceConnectionTokenRead(ctx, d, m)
}

func resourceConnectionTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", connectionTokenBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading connection token: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "account")
	setStringField(d, result, "connect_method")
	setStringField(d, result, "protocol")
	setObjectIDField(d, result, "asset")
	setStringField(d, result, "input_username")
	setStringField(d, result, "remote_addr")
	setBoolField(d, result, "is_active")
	setBoolField(d, result, "is_reusable")
	setStringField(d, result, "value")
	setStringField(d, result, "user_display")
	setStringField(d, result, "asset_display")
	setBoolField(d, result, "is_expired")
	setIntField(d, result, "expire_time")
	setStringField(d, result, "date_expired")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceConnectionTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"account":        d.Get("account").(string),
		"connect_method": d.Get("connect_method").(string),
	}
	if v, ok := d.GetOk("protocol"); ok {
		data["protocol"] = v.(string)
	}
	if v, ok := d.GetOk("asset"); ok {
		data["asset"] = v.(string)
	}
	if v, ok := d.GetOk("input_username"); ok {
		data["input_username"] = v.(string)
	}
	if v, ok := d.GetOk("input_secret"); ok {
		data["input_secret"] = v.(string)
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	data["is_active"] = d.Get("is_active").(bool)
	data["is_reusable"] = d.Get("is_reusable").(bool)

	resp, err := c.doRequest("PUT", connectionTokenBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating connection token: %s", resp.Status)
	}

	return resourceConnectionTokenRead(ctx, d, m)
}

func resourceConnectionTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", connectionTokenBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting connection token: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
