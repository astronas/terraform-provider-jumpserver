package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAssetCategory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetCategoryRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the asset category to look up.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the category.",
			},
		},
	}
}

func dataSourceAssetCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("assets/categories/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "value")

	return nil
}
