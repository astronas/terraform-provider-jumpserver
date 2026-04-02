package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ticketCommentBasePath = "tickets/comments/"

func resourceTicketComment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTicketCommentCreate,
		ReadContext:   resourceTicketCommentRead,
		UpdateContext: resourceTicketCommentUpdate,
		DeleteContext: resourceTicketCommentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"body": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The comment text.",
			},
			"user_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the comment author.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date.",
			},
		},
	}
}

func resourceTicketCommentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"body": d.Get("body").(string),
	}

	resp, err := c.doRequest("POST", ticketCommentBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating ticket comment: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving comment ID after creation, response: %v", result)
	}

	return resourceTicketCommentRead(ctx, d, m)
}

func resourceTicketCommentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", ticketCommentBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading ticket comment: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "body")
	setStringField(d, result, "user_display")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceTicketCommentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"body": d.Get("body").(string),
	}

	resp, err := c.doRequest("PUT", ticketCommentBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating ticket comment: %s", resp.Status)
	}

	return resourceTicketCommentRead(ctx, d, m)
}

func resourceTicketCommentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", ticketCommentBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting ticket comment: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
