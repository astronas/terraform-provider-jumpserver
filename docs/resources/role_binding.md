---
page_title: "jumpserver_role_binding Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages a role binding in JumpServer.
---

# jumpserver_role_binding

Manages a role binding in JumpServer. Role bindings assign a role to a user, optionally scoped to an organization.

## Example Usage

```hcl
resource "jumpserver_role_binding" "admin_binding" {
  user = "user-uuid-here"
  role = "role-uuid-here"
  org  = "org-uuid-here"
}
```

## Argument Reference

* `user` - (Required) The UUID of the user to bind the role to.
* `role` - (Required) The UUID of the role to bind.
* `org` - (Optional) The UUID of the organization. Required for org-scoped role bindings.

## Attribute Reference

* `id` - The UUID of the role binding.
* `scope` - The scope of the role binding (system or org).
* `org_name` - The name of the organization.
