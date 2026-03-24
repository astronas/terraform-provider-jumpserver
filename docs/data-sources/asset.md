---
page_title: "jumpserver_asset Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up an asset in JumpServer by name.
---

# `jumpserver_asset` Data Source

Use this data source to look up an existing asset in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_asset" "web_server" {
  name = "web-server-01"
}
```

## Argument Reference

* `name` - (Required) The name of the asset to look up.

## Attribute Reference

* `id` - The UUID of the asset.
* `address` - The address (IP or hostname) of the asset.
* `comment` - The description of the asset.
* `is_active` - Whether the asset is active.
