---
page_title: "jumpserver_org_role Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages an organization role in JumpServer.
---

# jumpserver_org_role

Manages an organization role in JumpServer. Organization roles define permissions scoped to an organization.

## Example Usage

```hcl
resource "jumpserver_org_role" "auditor" {
  name        = "custom-auditor"
  permissions = [1, 2, 3]
  comment     = "Custom auditor role"
}
```

## Argument Reference

* `name` - (Required) The name of the organization role.
* `permissions` - (Optional) List of permission IDs to assign to this role.
* `comment` - (Optional) A description for the role.

## Attribute Reference

* `id` - The UUID of the organization role.
* `builtin` - Whether this is a built-in role.
* `display_name` - The display name of the role.
* `users_amount` - The number of users assigned to this role.
