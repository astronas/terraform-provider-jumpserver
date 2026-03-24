package jumpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	bearerPrefix      = "Bearer "
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

// doRequest creates and executes an authenticated HTTP request.
// If body is not nil, it will be marshaled to JSON and Content-Type will be set.
func (c *Config) doRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.GetAPIEndpoint(path)

	if body != nil {
		jsonValue, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			return nil, err
		}
		req.Header.Set(headerContentType, contentTypeJSON)
		req.Header.Set("Authorization", bearerPrefix+c.Token)
		return c.NewHTTPClient().Do(req)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerPrefix+c.Token)
	return c.NewHTTPClient().Do(req)
}

// readEnumValue reads a field that may be returned as either a string
// or an object with {value, label} structure.
func readEnumValue(data map[string]interface{}, key string) (string, bool) {
	switch v := data[key].(type) {
	case map[string]interface{}:
		if val, ok := v["value"].(string); ok {
			return val, true
		}
	case string:
		return v, true
	}
	return "", false
}

// readObjectID reads a field that may be returned as either a string UUID
// or an object with an "id" field.
func readObjectID(data map[string]interface{}, key string) (string, bool) {
	switch v := data[key].(type) {
	case map[string]interface{}:
		if id, ok := v["id"].(string); ok {
			return id, true
		}
	case string:
		return v, true
	}
	return "", false
}

// readObjectIDs extracts IDs from an array that may contain strings or objects with an "id" field.
func readObjectIDs(items []interface{}) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		switch v := item.(type) {
		case string:
			ids = append(ids, v)
		case map[string]interface{}:
			if id, ok := v["id"].(string); ok {
				ids = append(ids, id)
			} else if id, ok := v["id"].(float64); ok {
				ids = append(ids, fmt.Sprintf("%d", int(id)))
			}
		}
	}
	return ids
}

// setStringField sets a string field on the ResourceData from the API response.
func setStringField(d *schema.ResourceData, data map[string]interface{}, key string) {
	if v, ok := data[key].(string); ok {
		d.Set(key, v)
	}
}

// setBoolField sets a bool field on the ResourceData from the API response.
func setBoolField(d *schema.ResourceData, data map[string]interface{}, key string) {
	if v, ok := data[key].(bool); ok {
		d.Set(key, v)
	}
}

// setIntField sets an int field on the ResourceData from the API response.
// JSON numbers are decoded as float64 by Go's encoding/json package.
func setIntField(d *schema.ResourceData, data map[string]interface{}, key string) {
	if v, ok := data[key].(float64); ok {
		d.Set(key, int(v))
	}
}

// setEnumField sets a field from the API response that may be a string or {value, label} object.
func setEnumField(d *schema.ResourceData, data map[string]interface{}, key string) {
	if v, ok := readEnumValue(data, key); ok {
		d.Set(key, v)
	}
}

// setObjectIDField sets a field from the API response that may be a string UUID or an object with id.
func setObjectIDField(d *schema.ResourceData, data map[string]interface{}, key string) {
	if v, ok := readObjectID(data, key); ok {
		d.Set(key, v)
	}
}

// setObjectIDsField sets a list field by extracting IDs from the API response array.
func setObjectIDsField(d *schema.ResourceData, data map[string]interface{}, key string) {
	if v, ok := data[key].([]interface{}); ok {
		d.Set(key, readObjectIDs(v))
	}
}
