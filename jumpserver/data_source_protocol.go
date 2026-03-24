package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProtocol() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProtocolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the protocol to look up.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The default port of the protocol.",
			},
		},
	}
}

func dataSourceProtocolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("assets/protocols/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setIntField(d, result, "port")

	return nil
}
