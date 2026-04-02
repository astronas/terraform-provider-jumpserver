package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const virtualAccountBasePath = "accounts/virtual-accounts/"

func resourceVirtualAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualAccountCreate,
		ReadContext:   resourceVirtualAccountRead,
		UpdateContext: resourceVirtualAccountUpdate,
		DeleteContext: resourceVirtualAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"@INPUT", "@USER", "@ANON", "@SPEC",
				}, false),
				Description: "The alias type of the virtual account.",
			},
			"secret_from_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to use the secret from the login account.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the virtual account.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the virtual account.",
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the virtual account.",
			},
		},
	}
}

func resourceVirtualAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"alias":             d.Get("alias").(string),
		"secret_from_login": d.Get("secret_from_login").(bool),
	}

	resp, err := c.doRequest("POST", virtualAccountBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_virtual_account", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating virtual account: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving virtual account ID after creation, response: %v", result)
	}

	return resourceVirtualAccountRead(ctx, d, m)
}

func resourceVirtualAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", virtualAccountBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading virtual account: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "alias")
	setBoolField(d, result, "secret_from_login")
	setStringField(d, result, "username")
	setStringField(d, result, "name")
	setStringField(d, result, "comment")

	return diags
}

func resourceVirtualAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"alias":             d.Get("alias").(string),
		"secret_from_login": d.Get("secret_from_login").(bool),
	}

	resp, err := c.doRequest("PUT", virtualAccountBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating virtual account: %s", resp.Status)
	}

	return resourceVirtualAccountRead(ctx, d, m)
}

func resourceVirtualAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", virtualAccountBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting virtual account: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
