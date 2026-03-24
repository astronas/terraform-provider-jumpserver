---
page_title: "jumpserver_favorite_asset Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a favorite asset in JumpServer.
---

# `jumpserver_favorite_asset` Resource

The `jumpserver_favorite_asset` resource allows you to mark an asset as a favorite in JumpServer. This is a create-only resource (no update).

## Example Usage

```hcl
resource "jumpserver_favorite_asset" "example" {
  asset = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`asset`** - (Required, ForceNew) The UUID of the asset to favorite.

## Attributes Reference

- **`id`** - The UUID of the favorite asset record.
