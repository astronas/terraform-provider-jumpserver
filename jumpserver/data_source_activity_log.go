package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceActivityLog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceActivityLogRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource ID to look up activity logs for.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The activity timestamp.",
			},
			"detail_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The detail URL for the activity.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The activity content.",
			},
			"r_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceActivityLogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	resourceID := d.Get("resource_id").(string)
	path := fmt.Sprintf("audits/activities/?resource_id=%s", resourceID)
	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error querying activity logs: %s", resp.Status)
	}

	var body interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return diag.FromErr(err)
	}

	var items []interface{}
	switch v := body.(type) {
	case map[string]interface{}:
		if results, ok := v["results"].([]interface{}); ok {
			items = results
		}
	case []interface{}:
		items = v
	}

	if len(items) == 0 {
		return diag.Errorf("no activity logs found for resource_id=%s", resourceID)
	}

	result, ok := items[0].(map[string]interface{})
	if !ok {
		return diag.Errorf("unexpected response format from audits/activities/")
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	}
	setStringField(d, result, "timestamp")
	setStringField(d, result, "detail_url")
	setStringField(d, result, "content")
	setStringField(d, result, "r_type")

	return nil
}
