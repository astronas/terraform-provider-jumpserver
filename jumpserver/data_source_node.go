package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodeRead,

		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value (display path) of the node to look up.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the node.",
			},
			"full_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path value of the node.",
			},
		},
	}
}

func dataSourceNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	value := d.Get("value").(string)
	result, err := c.dataSourceLookup("assets/nodes/", "value", value)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "name")
	setStringField(d, result, "full_value")

	return nil
}
