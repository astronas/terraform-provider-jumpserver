package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOperateLog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperateLogRead,

		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username to look up operate logs for.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action performed.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource affected.",
			},
			"resource": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource affected.",
			},
			"remote_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote address.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization ID.",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization name.",
			},
			"datetime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the operation.",
			},
		},
	}
}

func dataSourceOperateLogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	user := d.Get("user").(string)
	result, err := c.dataSourceLookup("audits/operate-logs/", "user", user)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setEnumField(d, result, "action")
	setStringField(d, result, "resource_type")
	setStringField(d, result, "resource")
	setStringField(d, result, "remote_addr")
	setStringField(d, result, "org_id")
	setStringField(d, result, "org_name")
	setStringField(d, result, "datetime")

	return nil
}
