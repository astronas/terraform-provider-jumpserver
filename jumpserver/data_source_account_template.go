package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account template to look up.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username defined in the template.",
			},
			"secret_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of secret (password, ssh_key, etc.).",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the account template is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
		},
	}
}

func dataSourceAccountTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("accounts/account-templates/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "username")
	setStringField(d, result, "secret_type")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")

	return nil
}
