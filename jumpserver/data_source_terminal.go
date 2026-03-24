package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTerminal() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTerminalRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the terminal to look up.",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote address.",
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
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
		},
	}
}

func dataSourceTerminalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("terminal/terminals/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "remote_addr")
	setEnumField(d, result, "type")
	setBoolField(d, result, "is_active")
	setBoolField(d, result, "is_alive")
	setStringField(d, result, "comment")

	return nil
}
