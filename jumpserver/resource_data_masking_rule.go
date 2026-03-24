package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const dataMaskingRuleBasePath = "acls/data-masking-rules/"

func resourceDataMaskingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataMaskingRuleCreate,
		ReadContext:   resourceDataMaskingRuleRead,
		UpdateContext: resourceDataMaskingRuleUpdate,
		DeleteContext: resourceDataMaskingRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the data masking rule.",
			},
			"users": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User filter as JSON (e.g. {\"type\":\"ids\",\"ids\":[\"uuid1\"]}).",
			},
			"assets": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Asset filter as JSON (e.g. {\"type\":\"ids\",\"ids\":[\"uuid1\"]}).",
			},
			"accounts": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of account usernames (e.g. [\"root\",\"@ALL\"]).",
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validation.IntBetween(1, 100),
				Description:  "Priority (1-100). Lower values match first.",
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "reject",
				ValidateFunc: validation.StringInSlice([]string{
					"reject", "accept", "review",
				}, false),
				Description: "Action when matched: reject, accept, review.",
			},
			"reviewers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of reviewer user IDs (required when action is review).",
			},
			"fields_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "password",
				Description: "Regex pattern matching fields to mask.",
			},
			"masking_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fixed_char",
				ValidateFunc: validation.StringInSlice([]string{
					"fixed_char",
				}, false),
				Description: "Masking method: fixed_char.",
			},
			"mask_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "######",
				Description: "Pattern used to mask the matched fields.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the rule is active.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the rule.",
			},
		},
	}
}

func resourceDataMaskingRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}
	var assets interface{}
	if err := json.Unmarshal([]byte(d.Get("assets").(string)), &assets); err != nil {
		return diag.Errorf("assets must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"users":          users,
		"assets":         assets,
		"accounts":       d.Get("accounts").([]interface{}),
		"priority":       d.Get("priority").(int),
		"action":         d.Get("action").(string),
		"fields_pattern": d.Get("fields_pattern").(string),
		"masking_method": d.Get("masking_method").(string),
		"mask_pattern":   d.Get("mask_pattern").(string),
		"is_active":      d.Get("is_active").(bool),
	}

	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}
	if v, ok := d.GetOk("comment"); ok {
		data["comment"] = v.(string)
	}

	resp, err := c.doRequest("POST", dataMaskingRuleBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating data masking rule: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving data masking rule ID after creation, response: %v", result)
	}

	return resourceDataMaskingRuleRead(ctx, d, m)
}

func resourceDataMaskingRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", dataMaskingRuleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading data masking rule: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setIntField(d, result, "priority")
	setBoolField(d, result, "is_active")
	setStringField(d, result, "comment")
	setStringField(d, result, "fields_pattern")
	setStringField(d, result, "mask_pattern")

	if v, ok := result["users"]; ok {
		b, _ := json.Marshal(v)
		d.Set("users", string(b))
	}
	if v, ok := result["assets"]; ok {
		b, _ := json.Marshal(v)
		d.Set("assets", string(b))
	}
	if v, ok := result["accounts"].([]interface{}); ok {
		d.Set("accounts", v)
	}

	setEnumField(d, result, "masking_method")
	setEnumField(d, result, "action")
	setObjectIDsField(d, result, "reviewers")

	return diags
}

func resourceDataMaskingRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var users interface{}
	if err := json.Unmarshal([]byte(d.Get("users").(string)), &users); err != nil {
		return diag.Errorf("users must be valid JSON: %s", err)
	}
	var assets interface{}
	if err := json.Unmarshal([]byte(d.Get("assets").(string)), &assets); err != nil {
		return diag.Errorf("assets must be valid JSON: %s", err)
	}

	data := map[string]interface{}{
		"name":           d.Get("name").(string),
		"users":          users,
		"assets":         assets,
		"accounts":       d.Get("accounts").([]interface{}),
		"priority":       d.Get("priority").(int),
		"action":         d.Get("action").(string),
		"fields_pattern": d.Get("fields_pattern").(string),
		"masking_method": d.Get("masking_method").(string),
		"mask_pattern":   d.Get("mask_pattern").(string),
		"is_active":      d.Get("is_active").(bool),
		"comment":        d.Get("comment").(string),
	}

	if v, ok := d.GetOk("reviewers"); ok {
		data["reviewers"] = v.([]interface{})
	}

	resp, err := c.doRequest("PUT", dataMaskingRuleBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating data masking rule: %s", resp.Status)
	}

	return resourceDataMaskingRuleRead(ctx, d, m)
}

func resourceDataMaskingRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", dataMaskingRuleBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting data masking rule: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
