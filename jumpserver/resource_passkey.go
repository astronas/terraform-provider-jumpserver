package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const passkeyBasePath = "authentication/passkeys/"

func resourcePasskey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePasskeyCreate,
		ReadContext:   resourcePasskeyRead,
		UpdateContext: resourcePasskeyUpdate,
		DeleteContext: resourcePasskeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the passkey.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the passkey is active.",
			},
			"platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The platform of the passkey.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the passkey.",
			},
			"date_last_used": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the passkey was last used.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
		},
	}
}

func resourcePasskeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"is_active": d.Get("is_active").(bool),
	}

	resp, err := c.doRequest("POST", passkeyBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating passkey: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving passkey ID after creation, response: %v", result)
	}

	return resourcePasskeyRead(ctx, d, m)
}

func resourcePasskeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", passkeyBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading passkey: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "platform")
	setStringField(d, result, "created_by")
	setStringField(d, result, "date_last_used")
	setStringField(d, result, "date_created")

	return diags
}

func resourcePasskeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"is_active": d.Get("is_active").(bool),
	}

	resp, err := c.doRequest("PUT", passkeyBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating passkey: %s", resp.Status)
	}

	return resourcePasskeyRead(ctx, d, m)
}

func resourcePasskeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", passkeyBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting passkey: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
