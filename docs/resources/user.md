---
page_title: "jumpserver_user Resource - terraform-provider-jumpserver"
subcategory: "Users & Groups"
description: |-
  Manages a user in JumpServer.
---

# `jumpserver_user` Resource

The jumpserver_user resource allows you to create and manage users in Jumpserver. Users represent individual accounts that can log into Jumpserver and interact with its resources.

## Example Usage

```hcl
resource "jumpserver_user" "example" {
  name         = "User 1"
  username     = "user1"
  email        = "user1@example.com"
  system_roles = ["cff4c2f3-mkj3-986j-b96f-31cf5hya1993"]
  is_active    = true
}
```

## Argument Reference

* `name` - (Required) The name of the user.
* `username` - (Required) The username of the user.
* `email` - (Required) The email of the user.
* `system_roles` - (Required) List of system roles assigned to the user.
* `is_active` - (Optional) Whether the user is active.

## Attribute Reference

* `id` - The ID of the user.
* `name` - The name of the user.
* `username` - The username of the user.
* `email` - The email of the user.
* `is_active` - Whether the user is active.
* `system_roles` - List of system roles assigned to the user.