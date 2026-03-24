---
page_title: "jumpserver_user_group Resource - terraform-provider-jumpserver"
subcategory: "Users & Groups"
description: |-
  Manages a user group in JumpServer.
---

# `jumpserver_user_group` Resource

The jumpserver_user_group resource allows you to create and manage user groups in JumpServer. User groups let you organize users and assign permissions collectively.

## Example Usage

```hcl
resource "jumpserver_user_group" "developers" {
  name    = "Developers"
  comment = "Development team"
  users   = [jumpserver_user.dev1.id, jumpserver_user.dev2.id]
}
```

## Argument Reference

* `name` - (Required) The name of the user group.
* `comment` - (Optional) A comment or description for the user group.
* `users` - (Optional) List of user IDs to include in this group.

## Attribute Reference

* `id` - The ID of the user group.
* `name` - The name of the user group.
* `comment` - The comment of the user group.
* `users` - List of user IDs in the group.
