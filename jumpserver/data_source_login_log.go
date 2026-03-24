package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLoginLog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoginLogRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username to look up login logs for.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The login type.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address.",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The city of origin.",
			},
			"user_agent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user agent string.",
			},
			"mfa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The MFA status.",
			},
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The login reason.",
			},
			"reason_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable login reason.",
			},
			"backend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication backend.",
			},
			"backend_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable backend name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The login status.",
			},
			"datetime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The login date and time.",
			},
		},
	}
}

func dataSourceLoginLogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	username := d.Get("username").(string)
	result, err := c.dataSourceLookup("audits/login-logs/", "username", username)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setEnumField(d, result, "type")
	setStringField(d, result, "ip")
	setStringField(d, result, "city")
	setStringField(d, result, "user_agent")
	setEnumField(d, result, "mfa")
	setStringField(d, result, "reason")
	setStringField(d, result, "reason_display")
	setStringField(d, result, "backend")
	setStringField(d, result, "backend_display")
	setEnumField(d, result, "status")
	setStringField(d, result, "datetime")

	return nil
}
