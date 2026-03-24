---
page_title: "jumpserver_system_role Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages a system role in JumpServer.
---

# jumpserver_system_role

Manages a system role in JumpServer. System roles define permissions at the system (global) level.

## Example Usage

```hcl
resource "jumpserver_system_role" "operator" {
  name        = "custom-operator"
  permissions = [10, 20, 30]
  comment     = "Custom operator role"
}
```

## Argument Reference

* `name` - (Required) The name of the system role.
* `permissions` - (Optional) List of permission IDs to assign to this role.
* `comment` - (Optional) A description for the role.

## Attribute Reference

* `id` - The UUID of the system role.
* `builtin` - Whether this is a built-in role.
* `display_name` - The display name of the role.
* `users_amount` - The number of users assigned to this role.
