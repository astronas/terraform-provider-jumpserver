package jumpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const chataiPromptBasePath = "settings/chatai-prompts/"

func resourceChatAIPrompt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChatAIPromptCreate,
		ReadContext:   resourceChatAIPromptRead,
		UpdateContext: resourceChatAIPromptUpdate,
		DeleteContext: resourceChatAIPromptDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ChatAI prompt.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content of the prompt.",
			},
			"builtin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the prompt is builtin.",
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

func resourceChatAIPromptCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":    d.Get("name").(string),
		"content": d.Get("content").(string),
		"builtin": d.Get("builtin").(bool),
	}

	resp, err := c.doRequest("POST", chataiPromptBasePath, data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating ChatAI prompt: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving ChatAI prompt ID after creation, response: %v", result)
	}

	return resourceChatAIPromptRead(ctx, d, m)
}

func resourceChatAIPromptRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", chataiPromptBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading ChatAI prompt: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, result, "name")
	setStringField(d, result, "content")
	setBoolField(d, result, "builtin")
	setStringField(d, result, "date_created")
	setStringField(d, result, "date_updated")

	return diags
}

func resourceChatAIPromptUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	data := map[string]interface{}{
		"name":    d.Get("name").(string),
		"content": d.Get("content").(string),
		"builtin": d.Get("builtin").(bool),
	}

	resp, err := c.doRequest("PUT", chataiPromptBasePath+d.Id()+"/", data)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating ChatAI prompt: %s", resp.Status)
	}

	return resourceChatAIPromptRead(ctx, d, m)
}

func resourceChatAIPromptDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	resp, err := c.doRequest("DELETE", chataiPromptBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting ChatAI prompt: %s", resp.Status)
	}

	d.SetId("")
	return diags
}
