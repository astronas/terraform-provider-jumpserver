---
page_title: "jumpserver_cloud Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a cloud asset in JumpServer by name.
---

# `jumpserver_cloud` Data Source

Use this data source to look up an existing cloud asset in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_cloud" "example" {
  name = "aws-instance-01"
}
```

## Argument Reference

- **`name`** - (Required) The name of the cloud asset to look up.

## Attributes Reference

- **`id`** - The UUID of the cloud asset.
- **`address`** - The address of the cloud asset.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the cloud asset is active.
