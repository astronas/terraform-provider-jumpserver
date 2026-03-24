---
page_title: "jumpserver_web Resource"
subcategory: "Assets"
description: |-
  Manages a web application asset in JumpServer.
---

# jumpserver_web

Manages a web application asset in JumpServer. Supports auto-login to web applications via CSS selectors for username, password, and submit button fields.

## Example Usage

```hcl
resource "jumpserver_web" "grafana" {
  name      = "grafana-prod"
  address   = "https://grafana.example.com"
  platform  = 10
  zone_name = "Production"
  node_name = "Web Apps"
  comment   = "Production Grafana"

  username_selector = "name=user"
  password_selector = "name=password"
  submit_selector   = "type=submit"

  accounts {
    name        = "admin"
    username    = "admin"
    secret_type = "password"
    secret      = "admin-password"
  }
}
```

## Argument Reference

### Base fields
- `name` - (Required) The name of the web asset.
- `address` - (Required) The URL of the web application.
- `platform` - (Required) The platform ID (integer).
- `zone_name` - (Required) The name of the zone.
- `node_name` - (Required) The name of the node.
- `comment` - (Optional) A description.
- `is_active` - (Optional) Whether the asset is active. Default: `true`.
- `accounts` - (Optional) List of account blocks.
- `protocols` - (Optional) List of protocol blocks (`name`, `port`).

### Web-specific fields
- `username_selector` - (Optional) CSS selector for username field. Default: `name=username`.
- `password_selector` - (Optional) CSS selector for password field. Default: `name=password`.
- `submit_selector` - (Optional) CSS selector for submit button. Default: `id=login_button`.
- `script` - (Optional) List of automation script steps for login.

## Attribute Reference

- `id` - The UUID of the web asset.
- `zone_id` - The resolved zone UUID.
- `node_ids` - The resolved node UUIDs.
