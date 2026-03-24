package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContentType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentTypeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the content type (resource type) to look up.",
			},
			"app_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The application label.",
			},
			"model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The model name.",
			},
		},
	}
}

func dataSourceContentTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("labels/resource-types/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "app_label")
	setStringField(d, result, "model")

	return nil
}
