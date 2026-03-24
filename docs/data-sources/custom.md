---
page_title: "jumpserver_custom Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a custom asset in JumpServer by name.
---

# `jumpserver_custom` Data Source

Use this data source to look up an existing custom asset in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_custom" "example" {
  name = "custom-device-01"
}
```

## Argument Reference

- **`name`** - (Required) The name of the custom asset to look up.

## Attributes Reference

- **`id`** - The UUID of the custom asset.
- **`address`** - The address of the custom asset.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the custom asset is active.
