package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const terminalRegistrationBasePath = "terminal/terminal-registrations/"

func resourceTerminalRegistration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTerminalRegistrationCreate,
		ReadContext:   resourceTerminalRegistrationRead,
		DeleteContext: resourceTerminalRegistrationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the terminal to register.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the terminal.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "A comment or description.",
			},
			"service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service account created for this terminal (JSON).",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote address.",
			},
		},
	}
}

func resourceTerminalRegistrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name": d.Get("name").(string),
		"type": d.Get("type").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", terminalRegistrationBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating terminal registration: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving terminal registration ID after creation, response: %v", result)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "remote_addr")
	if sa, ok := result["service_account"]; ok {
		saBytes, _ := json.Marshal(sa)
		d.Set("service_account", string(saBytes))
	}

	return nil
}

func resourceTerminalRegistrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceTerminalRegistrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceTerminalRegistrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
