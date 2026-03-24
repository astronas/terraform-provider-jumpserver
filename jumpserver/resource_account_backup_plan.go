package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const backupPlanBasePath = "accounts/account-backup-plans/"

func resourceAccountBackupPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountBackupPlanCreate,
		ReadContext:   resourceAccountBackupPlanRead,
		UpdateContext: resourceAccountBackupPlanUpdate,
		DeleteContext: resourceAccountBackupPlanDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the backup plan.",
			},
			"is_periodic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to run this backup periodically.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Execution interval in hours (1-65535).",
			},
			"crontab": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cron expression for scheduling (alternative to interval).",
			},
			"accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of account IDs to back up. Empty means all.",
			},
			"types": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Asset type filters.",
			},
			"nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of node IDs to include.",
			},
			"recipients_part_one": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs for email recipients (part one).",
			},
			"recipients_part_two": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs for email recipients (part two).",
			},
			"obj_recipients_part_one": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs for object storage recipients (part one).",
			},
			"obj_recipients_part_two": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs for object storage recipients (part two).",
			},
			"zip_encrypt_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password to encrypt the backup zip file.",
			},
			"is_password_divided_by_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to send password separately by email.",
			},
			"is_password_divided_by_obj_storage": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to send password separately to object storage.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the backup plan is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the backup plan.",
			},
		},
	}
}

func backupPlanBuildPayload(d *schema.ResourceData) map[string]interface{} {
	data := map[string]interface{}{
		"name":                                d.Get("name").(string),
		"is_periodic":                         d.Get("is_periodic").(bool),
		"interval":                            d.Get("interval").(int),
		"is_password_divided_by_email":        d.Get("is_password_divided_by_email").(bool),
		"is_password_divided_by_obj_storage":  d.Get("is_password_divided_by_obj_storage").(bool),
		"is_active":                           d.Get("is_active").(bool),
		"comment":                             d.Get("comment").(string),
	}

	if v, ok := d.GetOk("crontab"); ok {
		data["crontab"] = v.(string)
	}
	if v, ok := d.GetOk("accounts"); ok {
		data["accounts"] = v.([]interface{})
	}
	if v, ok := d.GetOk("types"); ok {
		data["types"] = v.([]interface{})
	}
	if v, ok := d.GetOk("nodes"); ok {
		data["nodes"] = v.([]interface{})
	}
	if v, ok := d.GetOk("recipients_part_one"); ok {
		data["recipients_part_one"] = v.([]interface{})
	}
	if v, ok := d.GetOk("recipients_part_two"); ok {
		data["recipients_part_two"] = v.([]interface{})
	}
	if v, ok := d.GetOk("obj_recipients_part_one"); ok {
		data["obj_recipients_part_one"] = v.([]interface{})
	}
	if v, ok := d.GetOk("obj_recipients_part_two"); ok {
		data["obj_recipients_part_two"] = v.([]interface{})
	}
	if v, ok := d.GetOk("zip_encrypt_password"); ok {
		data["zip_encrypt_password"] = v.(string)
	}

	return data
}

func resourceAccountBackupPlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := backupPlanBuildPayload(d)

	resp, err := c.doRequest("POST", backupPlanBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating account backup plan: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving account backup plan ID after creation, response: %v", result)
	}

	return resourceAccountBackupPlanRead(ctx, d, m)
}

func resourceAccountBackupPlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", backupPlanBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading account backup plan: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setBoolField(d, result, "is_periodic")
	setIntField(d, result, "interval")
	setStringField(d, result, "crontab")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")
	setBoolField(d, result, "is_password_divided_by_email")
	setBoolField(d, result, "is_password_divided_by_obj_storage")

	if v, ok := result["accounts"].([]interface{}); ok {
		d.Set("accounts", v)
	}
	if v, ok := result["types"].([]interface{}); ok {
		d.Set("types", v)
	}

	setObjectIDsField(d, result, "nodes")
	setObjectIDsField(d, result, "recipients_part_one")
	setObjectIDsField(d, result, "recipients_part_two")
	setObjectIDsField(d, result, "obj_recipients_part_one")
	setObjectIDsField(d, result, "obj_recipients_part_two")

	return diags
}

func resourceAccountBackupPlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := backupPlanBuildPayload(d)

	resp, err := c.doRequest("PUT", backupPlanBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating account backup plan: %s", resp.Status)
	}

	return resourceAccountBackupPlanRead(ctx, d, m)
}

func resourceAccountBackupPlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", backupPlanBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting account backup plan: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
