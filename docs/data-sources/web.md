---
page_title: "jumpserver_web Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a web asset in JumpServer by name.
---

# `jumpserver_web` Data Source

Use this data source to look up an existing web asset in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_web" "example" {
  name = "webapp-01"
}
```

## Argument Reference

- **`name`** - (Required) The name of the web asset to look up.

## Attributes Reference

- **`id`** - The UUID of the web asset.
- **`address`** - The address (URL) of the web asset.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the web asset is active.
