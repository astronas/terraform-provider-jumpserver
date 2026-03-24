---
page_title: "jumpserver_org_role_binding Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages an organization role binding in JumpServer.
---

# `jumpserver_org_role_binding` Resource

The `jumpserver_org_role_binding` resource allows you to assign organization-level roles to users within a specific organization.

## Example Usage

```hcl
resource "jumpserver_org_role_binding" "admin_binding" {
  user = "550e8400-e29b-41d4-a716-446655440000"
  role = "660e8400-e29b-41d4-a716-446655440001"
  org  = "770e8400-e29b-41d4-a716-446655440002"
}
```

## Argument Reference

* `user` - (Required) The UUID of the user.
* `role` - (Required) The UUID of the organization role.
* `org` - (Required) The UUID of the organization.

## Attribute Reference

* `id` - The UUID of the role binding.
* `scope` - The scope of the binding.
* `org_name` - The name of the organization.
