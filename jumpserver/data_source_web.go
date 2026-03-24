package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWeb() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the web asset to look up.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address (URL) of the web asset.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the web asset is active.",
			},
		},
	}
}

func dataSourceWebRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("assets/webs/", "name", name)
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
