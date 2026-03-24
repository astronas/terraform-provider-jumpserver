---
page_title: "jumpserver_endpoint Data Source - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Look up an endpoint by name in JumpServer.
---

# `jumpserver_endpoint` Data Source

Use this data source to look up an existing endpoint by name.

## Example Usage

```hcl
data "jumpserver_endpoint" "main" {
  name = "main-endpoint"
}
```

## Argument Reference

- **`name`** - (Required) The name of the endpoint to look up.

## Attribute Reference

- **`id`** - The UUID of the endpoint.
- **`host`** - The host address of the endpoint.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the endpoint is active.
