package jumpserver

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSiteMessage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiteMessageRead,

		Schema: map[string]*schema.Schema{
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subject of the site message to look up.",
			},
			"has_read": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the message has been read.",
			},
			"read_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the message was read.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message content (JSON).",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
		},
	}
}

func dataSourceSiteMessageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	subject := d.Get("subject").(string)
	result, err := c.dataSourceLookup("notifications/site-messages/", "subject", subject)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setBoolField(d, result, "has_read")
	setStringField(d, result, "read_at")
	setStringField(d, result, "date_created")
	if ct, ok := result["content"]; ok {
		ctBytes, _ := json.Marshal(ct)
		d.Set("content", string(ctBytes))
	}

	return nil
}
