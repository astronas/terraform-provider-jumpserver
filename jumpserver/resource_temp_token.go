package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const tempTokenBasePath = "authentication/temp-tokens/"

func resourceTempToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTempTokenCreate,
		ReadContext:   resourceTempTokenRead,
		DeleteContext: resourceTempTokenDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username associated with the token.",
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The token secret.",
			},
			"verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the token has been verified.",
			},
			"is_valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the token is valid.",
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
			"date_verified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The verification date.",
			},
			"date_expired": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date.",
			},
		},
	}
}

func resourceTempTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	resp, err := c.doRequest("POST", tempTokenBasePath, map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating temp token: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving temp token ID after creation, response: %v", result)
	}

	return resourceTempTokenRead(ctx, d, m)
}

func resourceTempTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", tempTokenBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading temp token: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "username")
	setStringField(d, result, "secret")
	setBoolField(d, result, "verified")
	setBoolField(d, result, "is_valid")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")
	setStringField(d, result, "date_verified")
	setStringField(d, result, "date_expired")

	return diags
}

func resourceTempTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	resp, err := c.doRequest("PATCH", tempTokenBasePath+d.Id()+"/", map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating temp token: %s", resp.Status)
	}

	return resourceTempTokenRead(ctx, d, m)
}

func resourceTempTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
