package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAssetPermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetPermissionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the asset permission to look up.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the permission is active.",
			},
			"date_start": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start date of the permission.",
			},
			"date_expired": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date of the permission.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
		},
	}
}

func dataSourceAssetPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("perms/asset-permissions/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setBoolField(d, result, "is_active")
	setStringField(d, result, "date_start")
	setStringField(d, result, "date_expired")
	setStringField(d, result, "comment")

	return nil
}
