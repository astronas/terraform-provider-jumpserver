---
page_title: "jumpserver_asset_perm_asset_relation Resource - terraform-provider-jumpserver"
subcategory: "Permissions"
description: |-
  Manages an asset permission to asset relation in JumpServer.
---

# `jumpserver_asset_perm_asset_relation` Resource

The `jumpserver_asset_perm_asset_relation` resource allows you to link an asset to an asset permission in JumpServer.

## Example Usage

```hcl
resource "jumpserver_asset_perm_asset_relation" "example" {
  asset           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  assetpermission = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`asset`** - (Required) The UUID of the asset.
- **`assetpermission`** - (Required) The UUID of the asset permission.

## Attributes Reference

- **`id`** - The integer ID of the relation.
- **`asset_display`** - The display name of the asset.
- **`assetpermission_display`** - The display name of the asset permission.
