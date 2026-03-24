package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const terminalCommandBasePath = "terminal/commands/"

func resourceTerminalCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTerminalCommandCreate,
		ReadContext:   resourceTerminalCommandRead,
		UpdateContext: resourceTerminalCommandUpdate,
		DeleteContext: resourceTerminalCommandDelete,

		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username.",
			},
			"asset": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The asset name.",
			},
			"input": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The command input.",
			},
			"session": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The session UUID.",
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account name.",
			},
			"output": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The command output.",
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The timestamp of the command.",
			},
			"risk_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The risk level: 0 (Ordinary), 1 (Dangerous).",
			},
			"org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The organization ID.",
			},
			"timestamp_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The formatted timestamp display.",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote address.",
			},
		},
	}
}

func resourceTerminalCommandCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":      d.Get("user").(string),
		"asset":     d.Get("asset").(string),
		"input":     d.Get("input").(string),
		"session":   d.Get("session").(string),
		"account":   d.Get("account").(string),
		"output":    d.Get("output").(string),
		"timestamp": d.Get("timestamp").(int),
	}
	if v, ok := d.GetOk("risk_level"); ok {
		data["risk_level"] = v.(int)
	}
	if v, ok := d.GetOk("org_id"); ok {
		data["org_id"] = v.(string)
	}

	resp, err := c.doRequest("POST", terminalCommandBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating terminal command: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving terminal command ID after creation, response: %v", result)
	}

	return resourceTerminalCommandRead(ctx, d, m)
}

func resourceTerminalCommandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", terminalCommandBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading terminal command: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "user")
	setStringField(d, result, "asset")
	setStringField(d, result, "input")
	setStringField(d, result, "session")
	setStringField(d, result, "account")
	setStringField(d, result, "output")
	setIntField(d, result, "timestamp")
	setIntField(d, result, "risk_level")
	setStringField(d, result, "org_id")
	setStringField(d, result, "timestamp_display")
	setStringField(d, result, "remote_addr")

	return diags
}

func resourceTerminalCommandUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"user":      d.Get("user").(string),
		"asset":     d.Get("asset").(string),
		"input":     d.Get("input").(string),
		"session":   d.Get("session").(string),
		"account":   d.Get("account").(string),
		"output":    d.Get("output").(string),
		"timestamp": d.Get("timestamp").(int),
	}
	if v, ok := d.GetOk("risk_level"); ok {
		data["risk_level"] = v.(int)
	}
	if v, ok := d.GetOk("org_id"); ok {
		data["org_id"] = v.(string)
	}

	resp, err := c.doRequest("PUT", terminalCommandBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating terminal command: %s", resp.Status)
	}

	return resourceTerminalCommandRead(ctx, d, m)
}

func resourceTerminalCommandDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", terminalCommandBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting terminal command: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
