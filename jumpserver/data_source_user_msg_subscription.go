package jumpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserMsgSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserMsgSubscriptionRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user UUID to look up message subscriptions for.",
			},
			"receive_backends": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of notification backend names the user is subscribed to.",
			},
		},
	}
}

func dataSourceUserMsgSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	userID := d.Get("user_id").(string)
	path := fmt.Sprintf("notifications/user-msg-subscription/%s/", userID)
	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading user message subscription: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userID)
	if v, ok := result["receive_backends"].([]interface{}); ok {
		d.Set("receive_backends", v)
	}

	return nil
}
