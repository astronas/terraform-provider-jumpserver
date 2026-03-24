package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account to look up.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the account.",
			},
			"asset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the associated asset.",
			},
			"secret_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of secret (password, ssh_key, etc.).",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the account is active.",
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("accounts/accounts/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "username")
	setObjectIDField(d, result, "asset")
	setStringField(d, result, "secret_type")
	setBoolField(d, result, "is_active")

	return nil
}
