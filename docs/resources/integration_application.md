---
page_title: "jumpserver_integration_application Resource - terraform-provider-jumpserver"
subcategory: "Infrastructure"
description: |-
  Manages a JumpServer Integration Application.
---

# jumpserver_integration_application

Manages a JumpServer Integration Application. Integration applications allow external services to access account secrets via API.

## Example Usage

```hcl
resource "jumpserver_integration_application" "ci_cd" {
  name     = "ci-cd-pipeline"
  accounts = [jumpserver_account.deploy.id]
  ip_group = ["10.0.0.0/8"]
  is_active = true
  comment  = "CI/CD pipeline access to deploy accounts"
}

resource "jumpserver_integration_application" "monitoring" {
  name     = "monitoring-app"
  accounts = [
    jumpserver_account.monitor_db.id,
    jumpserver_account.monitor_ssh.id,
  ]
  ip_group  = ["*"]
  is_active = true
}
```

## Argument Reference

- `name` - (Required) The name of the integration application.
- `accounts` - (Required) List of account IDs associated with this application.
- `ip_group` - (Optional) List of allowed IP ranges. Defaults to `["*"]` (all).
- `is_active` - (Optional) Whether the application is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the integration application.

## Import

```shell
terraform import jumpserver_integration_application.example <uuid>
```
