package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const jobBasePath = "ops/jobs/"

func resourceJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobCreate,
		ReadContext:   resourceJobRead,
		UpdateContext: resourceJobUpdate,
		DeleteContext: resourceJobDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the job.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "adhoc",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"adhoc", "playbook", "upload_file",
				}, false),
				Description: "The type of job: adhoc, playbook, or upload_file.",
			},
			"module": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shell",
				ValidateFunc: validation.StringInSlice([]string{
					"shell", "winshell", "python", "raw",
				}, false),
				Description: "The module to use for adhoc jobs.",
			},
			"args": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Command or arguments to execute.",
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the playbook to run (for playbook type jobs).",
			},
			"assets": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of asset UUIDs to run the job on.",
			},
			"nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of node UUIDs to run the job on.",
			},
			"runas_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "skip",
				ValidateFunc: validation.StringInSlice([]string{
					"skip", "if_need",
				}, false),
				Description: "Run-as policy: skip or if_need.",
			},
			"runas": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "root",
				Description: "The user to run as on the target assets.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Timeout in seconds (-1 for no timeout).",
			},
			"chdir": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Working directory on the target.",
			},
			"is_periodic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the job runs periodically.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Interval in hours for periodic execution.",
			},
			"crontab": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Crontab expression for scheduling.",
			},
			"run_after_save": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to run the job immediately after saving.",
			},
			"use_parameter_define": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to use parameter definitions.",
			},
			"parameters_define": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "{}",
				Description: "JSON string of parameter definitions.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description for the job.",
			},
		},
	}
}

func jobBuildPayload(d *schema.ResourceData) map[string]interface{} {
	data := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"type":                 d.Get("type").(string),
		"module":               d.Get("module").(string),
		"args":                 d.Get("args").(string),
		"runas_policy":         d.Get("runas_policy").(string),
		"runas":                d.Get("runas").(string),
		"timeout":              d.Get("timeout").(int),
		"chdir":                d.Get("chdir").(string),
		"is_periodic":          d.Get("is_periodic").(bool),
		"interval":             d.Get("interval").(int),
		"crontab":              d.Get("crontab").(string),
		"run_after_save":       d.Get("run_after_save").(bool),
		"use_parameter_define": d.Get("use_parameter_define").(bool),
		"comment":              d.Get("comment").(string),
	}

	if v, ok := d.GetOk("playbook"); ok {
		data["playbook"] = v.(string)
	}
	if v, ok := d.GetOk("assets"); ok {
		data["assets"] = v.([]interface{})
	}
	if v, ok := d.GetOk("nodes"); ok {
		data["nodes"] = v.([]interface{})
	}

	if v := d.Get("parameters_define").(string); v != "" && v != "{}" {
		var paramsDef interface{}
		if err := json.Unmarshal([]byte(v), &paramsDef); err == nil {
			data["parameters_define"] = paramsDef
		}
	}

	return data
}

func resourceJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := jobBuildPayload(d)

	resp, err := c.doRequest("POST", jobBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_job", d.Get("name").(string)); diags != nil {
		return diags
	}

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating job: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving job ID after creation, response: %v", result)
	}

	return resourceJobRead(ctx, d, m)
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", jobBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading job: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setEnumField(d, result, "type")
	setStringField(d, result, "module")
	setStringField(d, result, "args")
	setObjectIDField(d, result, "playbook")
	setObjectIDsField(d, result, "assets")
	setObjectIDsField(d, result, "nodes")
	setEnumField(d, result, "runas_policy")
	setStringField(d, result, "runas")
	setIntField(d, result, "timeout")
	setStringField(d, result, "chdir")
	setBoolField(d, result, "is_periodic")
	setIntField(d, result, "interval")
	setStringField(d, result, "crontab")
	setBoolField(d, result, "run_after_save")
	setBoolField(d, result, "use_parameter_define")
	setStringField(d, result, "comment")

	// parameters_define: serialize back to JSON string
	if v, ok := result["parameters_define"]; ok && v != nil {
		jsonBytes, err := json.Marshal(v)
		if err == nil {
			d.Set("parameters_define", string(jsonBytes))
		}
	}

	return diags
}

func resourceJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := jobBuildPayload(d)

	resp, err := c.doRequest("PUT", jobBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating job: %s", resp.Status)
	}

	return resourceJobRead(ctx, d, m)
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", jobBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting job: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
