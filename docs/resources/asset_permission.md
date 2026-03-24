---
page_title: "jumpserver_asset_permission Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages an asset permission in JumpServer.
---

# `jumpserver_asset_permission` Resource

The jumpserver_asset_permission resource allows you to create and manage asset permissions in Jumpserver. Asset permissions define which users (or user groups) have access to which assets (or nodes) using which system users, and what actions are allowed.

## Example Usage

```hcl
resource "jumpserver_asset_permission" "example_permission" {
  name         = "ops-to-prod"
  is_active    = true
  users        = ["<user-uuid>"]
  user_groups  = ["<user-group-uuid>"]
  assets       = ["<asset-uuid>"]
  nodes        = ["<node-uuid>"]
  system_users = ["<system-user-uuid>"]
  actions      = ["connect"]
}
```

## Argument Reference

* `name` - (Required) The name of the asset permission rule.
* `is_active` - (Optional) Whether the permission rule is active. Defaults to `true`.
* `users` - (Optional) List of user IDs to grant permission to.
* `user_groups` - (Optional) List of user group IDs to grant permission to.
* `assets` - (Optional) List of asset IDs covered by this permission.
* `nodes` - (Optional) List of node IDs covered by this permission.
* `system_users` - (Optional) List of system user IDs allowed by this permission.
* `actions` - (Optional) List of allowed actions (e.g. `connect`, `upload`, `download`).

## Attribute Reference

* `id` - The ID of the asset permission rule.
