---
page_title: "jumpserver_asset_permission Data Source - terraform-provider-jumpserver"
subcategory: "Permissions"
description: |-
  Look up an asset permission in JumpServer by name.
---

# `jumpserver_asset_permission` Data Source

Use this data source to look up an existing asset permission in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_asset_permission" "example" {
  name = "dev-team-access"
}
```

## Argument Reference

- **`name`** - (Required) The name of the asset permission to look up.

## Attributes Reference

- **`id`** - The UUID of the asset permission.
- **`is_active`** - Whether the permission is active.
- **`date_start`** - The start date.
- **`date_expired`** - The expiration date.
- **`comment`** - A comment or description.
