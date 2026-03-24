---
page_title: "jumpserver_database Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a database asset by name in JumpServer.
---

# `jumpserver_database` Data Source

Use this data source to look up an existing database asset by name.

## Example Usage

```hcl
data "jumpserver_database" "prod_db" {
  name = "production-mysql"
}
```

## Argument Reference

- **`name`** - (Required) The name of the database asset to look up.

## Attribute Reference

- **`id`** - The UUID of the database asset.
- **`address`** - The address of the database.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the database asset is active.
