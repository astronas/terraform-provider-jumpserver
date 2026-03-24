---
page_title: "jumpserver_asset_category Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up an asset category in JumpServer by name.
---

# `jumpserver_asset_category` Data Source

Use this data source to look up an existing asset category in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_asset_category" "example" {
  name = "host"
}
```

## Argument Reference

- **`name`** - (Required) The name of the asset category to look up.

## Attributes Reference

- **`id`** - The ID of the category.
- **`value`** - The value of the category.
