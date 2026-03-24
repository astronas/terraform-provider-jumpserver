package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustom() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the custom asset to look up.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the custom asset.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the custom asset is active.",
			},
		},
	}
}

func dataSourceCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("assets/customs/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "address")
	setStringField(d, result, "comment")
	setBoolField(d, result, "is_active")

	return nil
}
