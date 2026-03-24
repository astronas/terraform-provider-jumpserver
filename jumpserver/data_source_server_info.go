package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServerInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerInfoRead,

		Schema: map[string]*schema.Schema{
			"current_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server's current time.",
			},
		},
	}
}

func dataSourceServerInfoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	resp, err := c.doRequest("GET", "settings/server-info/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading server info: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("server-info")
	if v, ok := result["CURRENT_TIME"].(string); ok {
		d.Set("current_time", v)
	}

	return nil
}
