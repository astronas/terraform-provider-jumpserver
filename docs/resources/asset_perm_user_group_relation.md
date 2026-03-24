---
page_title: "jumpserver_asset_perm_user_group_relation Resource - terraform-provider-jumpserver"
subcategory: "Permissions"
description: |-
  Manages an asset permission to user group relation in JumpServer.
---

# `jumpserver_asset_perm_user_group_relation` Resource

The `jumpserver_asset_perm_user_group_relation` resource allows you to link a user group to an asset permission in JumpServer.

## Example Usage

```hcl
resource "jumpserver_asset_perm_user_group_relation" "example" {
  usergroup       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  assetpermission = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`usergroup`** - (Required) The UUID of the user group.
- **`assetpermission`** - (Required) The UUID of the asset permission.

## Attributes Reference

- **`id`** - The integer ID of the relation.
- **`usergroup_display`** - The display name of the user group.
- **`assetpermission_display`** - The display name of the asset permission.
