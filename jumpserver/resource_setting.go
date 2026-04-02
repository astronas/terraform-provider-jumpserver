package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const settingBasePath = "settings/setting/"

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingCreate,
		ReadContext:   resourceSettingRead,
		UpdateContext: resourceSettingUpdate,
		DeleteContext: resourceSettingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"site_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The site URL (uri format).",
			},
			"user_guide_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The user guide URL.",
			},
			"global_org_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The global organization display name.",
			},
			"help_document_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The help document URL.",
			},
			"help_support_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The help support URL.",
			},
		},
	}
}

func resourceSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("settings")
	return resourceSettingUpdate(ctx, d, m)
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", settingBasePath, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading settings: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := result["SITE_URL"].(string); ok {
		d.Set("site_url", v)
	}
	if v, ok := result["USER_GUIDE_URL"].(string); ok {
		d.Set("user_guide_url", v)
	}
	if v, ok := result["GLOBAL_ORG_DISPLAY_NAME"].(string); ok {
		d.Set("global_org_display_name", v)
	}
	if v, ok := result["HELP_DOCUMENT_URL"].(string); ok {
		d.Set("help_document_url", v)
	}
	if v, ok := result["HELP_SUPPORT_URL"].(string); ok {
		d.Set("help_support_url", v)
	}

	return diags
}

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"SITE_URL": d.Get("site_url").(string),
	}
	if v, ok := d.GetOk("user_guide_url"); ok {
		data["USER_GUIDE_URL"] = v.(string)
	}
	if v, ok := d.GetOk("global_org_display_name"); ok {
		data["GLOBAL_ORG_DISPLAY_NAME"] = v.(string)
	}
	if v, ok := d.GetOk("help_document_url"); ok {
		data["HELP_DOCUMENT_URL"] = v.(string)
	}
	if v, ok := d.GetOk("help_support_url"); ok {
		data["HELP_SUPPORT_URL"] = v.(string)
	}

	resp, err := c.doRequest("PUT", settingBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating settings: %s", resp.Status)
	}

	return resourceSettingRead(ctx, d, m)
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
