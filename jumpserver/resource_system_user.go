package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const systemUserBasePath = "assets/system-users/"

func resourceSystemUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemUserCreate,
		ReadContext:   resourceSystemUserRead,
		UpdateContext: resourceSystemUserUpdate,
		DeleteContext: resourceSystemUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "common",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ssh",
			},
			"login_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "auto",
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  81,
			},
			"sudo": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/bin/whoami",
			},
			"shell": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/bin/bash",
			},
			"sftp_root": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tmp",
			},
			"home": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/home/student",
			},
			"username_same_with_user": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"su_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceSystemUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	// Prepare payload
	payload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"username":                d.Get("username").(string),
		"password":                d.Get("password").(string),
		"type":                    d.Get("type").(string),
		"protocol":                d.Get("protocol").(string),
		"login_mode":              d.Get("login_mode").(string),
		"priority":                d.Get("priority").(int),
		"sudo":                    d.Get("sudo").(string),
		"shell":                   d.Get("shell").(string),
		"sftp_root":               d.Get("sftp_root").(string),
		"home":                    d.Get("home").(string),
		"username_same_with_user": d.Get("username_same_with_user").(bool),
		"auto_push":               d.Get("auto_push").(bool),
		"su_enabled":              d.Get("su_enabled").(bool),
	}

	resp, err := c.doRequest("POST", systemUserBasePath, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_system_user", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Failed to create system user. Status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID
	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Failed to retrieve ID for created system user")
	}

	return diags
}

func resourceSystemUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()

	resp, err := c.doRequest("GET", systemUserBasePath+id+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error fetching system user: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "username")
	setStringField(d, result, "type")
	setStringField(d, result, "protocol")
	setStringField(d, result, "login_mode")

	// Handle numeric values with proper type conversion
	if priority, ok := result["priority"].(float64); ok {
		d.Set("priority", int(priority))
	} else {
		return diag.Errorf("Failed to parse 'priority' field from API response")
	}

	// Additional fields to set if available
	setStringField(d, result, "home")
	setStringField(d, result, "shell")

	d.SetId(id)

	return diags
}

func resourceSystemUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	// Prepare payload for update
	payload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"username":                d.Get("username").(string),
		"password":                d.Get("password").(string),
		"type":                    d.Get("type").(string),
		"protocol":                d.Get("protocol").(string),
		"login_mode":              d.Get("login_mode").(string),
		"priority":                d.Get("priority").(int),
		"sudo":                    d.Get("sudo").(string),
		"shell":                   d.Get("shell").(string),
		"sftp_root":               d.Get("sftp_root").(string),
		"home":                    d.Get("home").(string),
		"username_same_with_user": d.Get("username_same_with_user").(bool),
		"auto_push":               d.Get("auto_push").(bool),
		"su_enabled":              d.Get("su_enabled").(bool),
	}

	id := d.Id()

	resp, err := c.doRequest("PUT", systemUserBasePath+id+"/", payload)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to update system user. Status code: %d", resp.StatusCode)
	}

	return resourceSystemUserRead(ctx, d, m)
}

func resourceSystemUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()

	resp, err := c.doRequest("DELETE", systemUserBasePath+id+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Failed to delete system user. Status code: %d", resp.StatusCode)
	}

	d.SetId("") // Mark resource as deleted

	return diags
}
