---
page_title: "jumpserver_asset_perm_user_relation Resource - terraform-provider-jumpserver"
subcategory: "Permissions"
description: |-
  Manages an asset permission to user relation in JumpServer.
---

# `jumpserver_asset_perm_user_relation` Resource

The `jumpserver_asset_perm_user_relation` resource allows you to link a user to an asset permission in JumpServer.

## Example Usage

```hcl
resource "jumpserver_asset_perm_user_relation" "example" {
  user            = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  assetpermission = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`user`** - (Required) The UUID of the user.
- **`assetpermission`** - (Required) The UUID of the asset permission.

## Attributes Reference

- **`id`** - The integer ID of the relation.
- **`user_display`** - The display name of the user.
- **`assetpermission_display`** - The display name of the asset permission.
