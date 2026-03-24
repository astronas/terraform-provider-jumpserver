package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSmsBackend() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSmsBackendRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the SMS backend to look up.",
			},
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable label of the backend.",
			},
		},
	}
}

func dataSourceSmsBackendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("settings/sms/backend/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	setStringField(d, result, "label")

	return nil
}
