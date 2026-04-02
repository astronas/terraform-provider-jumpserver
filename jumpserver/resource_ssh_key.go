package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const sshKeyBasePath = "authentication/ssh-key/"

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceSSHKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the SSH key.",
			},
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The public key content.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the SSH key is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the SSH key.",
			},
			"public_key_comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Comment extracted from the public key.",
			},
			"public_key_hash_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MD5 hash of the public key.",
			},
			"date_last_used": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the SSH key was last used.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date.",
			},
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"is_active": d.Get("is_active").(bool),
	}
	if v, ok := d.GetOk("public_key"); ok {
		data["public_key"] = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", sshKeyBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating SSH key: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving SSH key ID after creation, response: %v", result)
	}

	return resourceSSHKeyRead(ctx, d, m)
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", sshKeyBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading SSH key: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "public_key")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")
	setStringField(d, result, "public_key_comment")
	setStringField(d, result, "public_key_hash_md5")
	setStringField(d, result, "date_last_used")
	setStringField(d, result, "date_created")

	return diags
}

func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":      d.Get("name").(string),
		"is_active": d.Get("is_active").(bool),
		"comment":   d.Get("comment").(string),
	}
	if v, ok := d.GetOk("public_key"); ok {
		data["public_key"] = v.(string)
	}

	resp, err := c.doRequest("PUT", sshKeyBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating SSH key: %s", resp.Status)
	}

	return resourceSSHKeyRead(ctx, d, m)
}

func resourceSSHKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", sshKeyBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting SSH key: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
