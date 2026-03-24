package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLabel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLabelRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the label to look up.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the label.",
			},
			"color": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The color of the label.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
		},
	}
}

func dataSourceLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("labels/labels/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "value")
	setStringField(d, result, "color")
	setStringField(d, result, "comment")

	return nil
}
