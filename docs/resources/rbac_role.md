---
page_title: "jumpserver_rbac_role Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages an RBAC role in JumpServer.
---

# `jumpserver_rbac_role` Resource

The `jumpserver_rbac_role` resource allows you to create and manage RBAC roles in JumpServer.

## Example Usage

```hcl
resource "jumpserver_rbac_role" "custom_role" {
  name        = "CustomOperator"
  comment     = "Custom operator role"
  permissions = [1, 2, 3]
}
```

## Argument Reference

- **`name`** - (Required) The name of the role (max 128 characters).
- **`permissions`** - (Optional) List of permission IDs to assign (write-only).
- **`comment`** - (Optional) A description for the role.

## Attributes Reference

- **`id`** - The UUID of the role.
- **`display_name`** - The display name of the role.
- **`scope`** - The scope of the role.
- **`builtin`** - Whether this is a built-in role.
- **`users_amount`** - The number of users assigned to this role.
- **`created_by`** - Who created this role.
- **`updated_by`** - Who last updated this role.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
