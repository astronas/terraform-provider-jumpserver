---
page_title: "jumpserver_gateway Resource"
subcategory: "Assets"
description: |-
  Manages a gateway asset in JumpServer.
---

# jumpserver_gateway

Manages a gateway asset in JumpServer. Gateways serve as network proxies for zones — when connecting to assets within a zone, the connection is routed through the gateway.

## Example Usage

```hcl
resource "jumpserver_gateway" "dmz_gw" {
  name      = "dmz-gateway"
  address   = "10.0.1.1"
  platform  = 1
  zone_name = "DMZ"
  node_name = "Gateways"
  comment   = "Gateway proxy for DMZ zone"

  protocols {
    name = "ssh"
    port = 22
  }

  accounts {
    name        = "gw-admin"
    username    = "admin"
    secret_type = "password"
    secret      = "gw-password"
  }
}
```

## Argument Reference

- `name` - (Required) The name of the gateway.
- `address` - (Required) IP address or hostname of the gateway.
- `platform` - (Required) The platform ID (integer).
- `zone_name` - (Required) The name of the zone this gateway serves.
- `node_name` - (Required) The name of the node.
- `comment` - (Optional) A description.
- `is_active` - (Optional) Whether the gateway is active. Default: `true`.
- `accounts` - (Optional) List of account blocks.
- `protocols` - (Optional) List of protocol blocks (`name`, `port`).

## Attribute Reference

- `id` - The UUID of the gateway asset.
- `zone_id` - The resolved zone UUID.
- `node_ids` - The resolved node UUIDs.
