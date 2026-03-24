package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrgRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrgRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization role to look up.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the role.",
			},
			"builtin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a built-in role.",
			},
		},
	}
}

func dataSourceOrgRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("rbac/org-roles/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "name")
	setStringField(d, result, "comment")
	setBoolField(d, result, "builtin")

	return nil
}
