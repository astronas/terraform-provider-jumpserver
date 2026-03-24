package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const terminalBasePath = "terminal/terminals/"

func resourceTerminal() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTerminalCreate,
		ReadContext:   resourceTerminalRead,
		UpdateContext: resourceTerminalUpdate,
		DeleteContext: resourceTerminalDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the terminal.",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The remote address of the terminal.",
			},
			"command_storage": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The command storage backend name.",
			},
			"replay_storage": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The replay storage backend name.",
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
				Description: "The type of the terminal.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the terminal is active.",
			},
			"is_alive": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the terminal is alive.",
			},
			"session_online": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of online sessions.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
		},
	}
}

func resourceTerminalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name": d.Get("name").(string),
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	if v, ok := d.GetOk("command_storage"); ok {
		data["command_storage"] = v.(string)
	}
	if v, ok := d.GetOk("replay_storage"); ok {
		data["replay_storage"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", terminalBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating terminal: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving terminal ID after creation, response: %v", result)
	}

	return resourceTerminalRead(ctx, d, m)
}

func resourceTerminalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", terminalBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading terminal: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "remote_addr")
	setStringField(d, result, "command_storage")
	setStringField(d, result, "replay_storage")
	setStringField(d, result, "comment")
	setEnumField(d, result, "type")
	setBoolField(d, result, "is_active")
	setBoolField(d, result, "is_alive")
	setIntField(d, result, "session_online")
	setStringField(d, result, "date_created")

	return diags
}

func resourceTerminalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name": d.Get("name").(string),
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		data["remote_addr"] = v.(string)
	}
	if v, ok := d.GetOk("command_storage"); ok {
		data["command_storage"] = v.(string)
	}
	if v, ok := d.GetOk("replay_storage"); ok {
		data["replay_storage"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("PUT", terminalBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating terminal: %s", resp.Status)
	}

	return resourceTerminalRead(ctx, d, m)
}

func resourceTerminalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", terminalBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting terminal: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
