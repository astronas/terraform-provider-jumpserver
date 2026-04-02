package jumpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const userBasePath = "users/users/"

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_roles": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	user := map[string]interface{}{
		"name":         d.Get("name").(string),
		"username":     d.Get("username").(string),
		"email":        d.Get("email").(string),
		"is_active":    d.Get("is_active").(bool),
		"system_roles": d.Get("system_roles").([]interface{}),
	}

	resp, err := c.doRequest("POST", userBasePath, user)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if diags := checkAlreadyExists(resp, "jumpserver_user", d.Get("name").(string)); diags != nil {
		return diags
	}

	// Check for 201 Created status code
	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating user: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Log the entire response
	log.Printf("Response Body: %v\n", result)

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving user ID after creation, response: %v", result)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	resp, err := c.doRequest("GET", userBasePath+d.Id()+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading user: %s", resp.Status)
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return diag.FromErr(err)
	}

	setStringField(d, user, "name")
	setStringField(d, user, "username")
	setStringField(d, user, "email")
	setBoolField(d, user, "is_active")
	d.Set("system_roles", user["system_roles"].([]interface{}))

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	user := map[string]interface{}{
		"name":         d.Get("name").(string),
		"username":     d.Get("username").(string),
		"email":        d.Get("email").(string),
		"is_active":    d.Get("is_active").(bool),
		"system_roles": d.Get("system_roles").([]interface{}),
	}

	resp, err := c.doRequest("PUT", userBasePath+d.Id()+"/", user)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating user: %s", resp.Status)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()

	resp, err := c.doRequest("DELETE", userBasePath+id+"/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for 204 No Content status code
	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting user: %s", resp.Status)
	}

	d.SetId("") // Mark resource as destroyed
	return diags
}
