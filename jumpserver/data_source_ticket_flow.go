package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTicketFlow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTicketFlowRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ticket flow to look up.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the ticket flow.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the ticket flow is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comment or description.",
			},
		},
	}
}

func dataSourceTicketFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("tickets/flows/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setEnumField(d, result, "type")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")

	return nil
}
