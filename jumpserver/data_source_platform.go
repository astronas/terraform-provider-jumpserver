package jumpserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlatform() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlatformRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the platform to look up.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the platform.",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The category of the platform.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the platform.",
			},
		},
	}
}

func dataSourcePlatformRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("assets/platforms/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(float64); ok {
		d.SetId(fmt.Sprintf("%d", int(id)))
	} else if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "name")
	setStringField(d, result, "comment")
	setEnumField(d, result, "category")
	setEnumField(d, result, "type")

	return nil
}
