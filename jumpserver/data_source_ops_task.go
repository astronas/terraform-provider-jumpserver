package jumpserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpsTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpsTaskRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ops task to look up.",
			},
			"meta": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Task metadata.",
			},
		},
	}
}

func dataSourceOpsTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	name := d.Get("name").(string)
	result, err := c.dataSourceLookup("ops/tasks/", "name", name)
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	if meta, ok := result["meta"].(map[string]interface{}); ok {
		flatMeta := make(map[string]string)
		for k, v := range meta {
			if s, ok := v.(string); ok {
				flatMeta[k] = s
			}
		}
		d.Set("meta", flatMeta)
	}

	return nil
}
