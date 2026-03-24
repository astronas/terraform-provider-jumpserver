package jumpserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRbacPermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRbacPermissionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the permission to look up.",
			},
			"codename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The codename of the permission.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type ID.",
			},
		},
	}
}

func dataSourceRbacPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("rbac/permissions/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(float64); ok {
		d.SetId(fmt.Sprintf("%d", int(id)))
	}
	setStringField(d, result, "codename")
	setObjectIDField(d, result, "content_type")

	return nil
}
