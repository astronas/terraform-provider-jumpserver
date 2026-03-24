---
page_title: "jumpserver_system_role_binding Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages a system role binding in JumpServer.
---

# `jumpserver_system_role_binding` Resource

The `jumpserver_system_role_binding` resource allows you to assign system-level roles to users.

## Example Usage

```hcl
resource "jumpserver_system_role_binding" "sysadmin" {
  user = "550e8400-e29b-41d4-a716-446655440000"
  role = "660e8400-e29b-41d4-a716-446655440001"
}
```

## Argument Reference

* `user` - (Required) The UUID of the user.
* `role` - (Required) The UUID of the system role.

## Attribute Reference

* `id` - The UUID of the role binding.
* `scope` - The scope of the binding.
* `org_name` - The name of the organization.
