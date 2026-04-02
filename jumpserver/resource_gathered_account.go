package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const gatheredAccountBasePath = "accounts/gathered-accounts/"

func resourceGatheredAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGatheredAccountCreate,
		ReadContext:   resourceGatheredAccountRead,
		UpdateContext: resourceGatheredAccountUpdate,
		DeleteContext: resourceGatheredAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"asset": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the asset.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The discovered account username.",
			},
			"address_last_login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address of last login.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the gathered account.",
			},
			"remote_present": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the account is remotely present.",
			},
			"present": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the account is present.",
			},
			"date_last_login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date of last login.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date.",
			},
		},
	}
}

func resourceGatheredAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"asset": d.Get("asset").(string),
	}

	resp, err := c.doRequest("POST", gatheredAccountBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_gathered_account", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating gathered account: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving gathered account ID after creation, response: %v", result)
	}

	return resourceGatheredAccountRead(ctx, d, m)
}

func resourceGatheredAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", gatheredAccountBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading gathered account: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setObjectIDField(d, result, "asset")
	setStringField(d, result, "username")
	setStringField(d, result, "address_last_login")
	setEnumField(d, result, "status")
	setBoolField(d, result, "remote_present")
	setBoolField(d, result, "present")
	setStringField(d, result, "date_last_login")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceGatheredAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"asset": d.Get("asset").(string),
	}

	resp, err := c.doRequest("PUT", gatheredAccountBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating gathered account: %s", resp.Status)
	}

	return resourceGatheredAccountRead(ctx, d, m)
}

func resourceGatheredAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", gatheredAccountBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting gathered account: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
