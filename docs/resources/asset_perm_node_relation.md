---
page_title: "jumpserver_asset_perm_node_relation Resource - terraform-provider-jumpserver"
subcategory: "Permissions"
description: |-
  Manages an asset permission to node relation in JumpServer.
---

# `jumpserver_asset_perm_node_relation` Resource

The `jumpserver_asset_perm_node_relation` resource allows you to link a node to an asset permission in JumpServer.

## Example Usage

```hcl
resource "jumpserver_asset_perm_node_relation" "example" {
  node            = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  assetpermission = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`node`** - (Required) The UUID of the node.
- **`assetpermission`** - (Required) The UUID of the asset permission.

## Attributes Reference

- **`id`** - The integer ID of the relation.
- **`node_display`** - The display name of the node.
- **`assetpermission_display`** - The display name of the asset permission.
