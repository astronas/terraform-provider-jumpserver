---
page_title: "jumpserver_custom Resource"
subcategory: "Assets"
description: |-
  Manages a custom asset in JumpServer.
---

# jumpserver_custom

Manages a custom asset in JumpServer. Used for asset types not covered by the standard types (host, database, device, web, cloud, gateway). Supports arbitrary configuration via a JSON field.

## Example Usage

```hcl
resource "jumpserver_custom" "internal_tool" {
  name      = "internal-tool"
  address   = "tool.internal.example.com"
  platform  = 30
  zone_name = "Production"
  node_name = "Custom"
  comment   = "Internal custom tool"

  custom_info = jsonencode({
    tool_type = "monitoring"
    version   = "2.1"
  })

  protocols {
    name = "http"
    port = 8080
  }
}
```

## Argument Reference

### Base fields
- `name` - (Required) The name of the custom asset.
- `address` - (Required) The address of the asset.
- `platform` - (Required) The platform ID (integer).
- `zone_name` - (Required) The name of the zone.
- `node_name` - (Required) The name of the node.
- `comment` - (Optional) A description.
- `is_active` - (Optional) Whether the asset is active. Default: `true`.
- `accounts` - (Optional) List of account blocks.
- `protocols` - (Optional) List of protocol blocks (`name`, `port`).

### Custom-specific fields
- `custom_info` - (Optional) Custom configuration as a JSON string. Default: `{}`.

## Attribute Reference

- `id` - The UUID of the custom asset.
- `zone_id` - The resolved zone UUID.
- `node_ids` - The resolved node UUIDs.
