---
page_title: "jumpserver_gateway Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a gateway by name in JumpServer.
---

# `jumpserver_gateway` Data Source

Use this data source to look up an existing gateway by name.

## Example Usage

```hcl
data "jumpserver_gateway" "main_gw" {
  name = "main-gateway"
}
```

## Argument Reference

- **`name`** - (Required) The name of the gateway to look up.

## Attribute Reference

- **`id`** - The UUID of the gateway.
- **`address`** - The address of the gateway.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the gateway is active.
