---
page_title: "jumpserver_user_group_relation Resource - terraform-provider-jumpserver"
subcategory: "Users"
description: |-
  Manages a user-to-group relationship in JumpServer.
---

# `jumpserver_user_group_relation` Resource

Manages an explicit relationship between a user and a user group in JumpServer.

## Example Usage

```hcl
resource "jumpserver_user_group_relation" "example" {
  user      = jumpserver_user.example.id
  usergroup = jumpserver_user_group.example.id
}
```

## Argument Reference

- **`user`** - (Required) The UUID of the user.
- **`usergroup`** - (Required) The UUID of the user group.

## Attribute Reference

- **`id`** - The ID of the relation.
- **`user_display`** - The display name of the user.
- **`usergroup_display`** - The display name of the user group.
