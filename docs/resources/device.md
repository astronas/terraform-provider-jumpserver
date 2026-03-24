---
page_title: "jumpserver_device Resource"
subcategory: "Assets"
description: |-
  Manages a device asset in JumpServer.
---

# jumpserver_device

Manages a device asset in JumpServer. Used for generic infrastructure devices such as network equipment, servers, firewalls, etc.

## Example Usage

```hcl
resource "jumpserver_device" "switch_core" {
  name      = "core-switch-01"
  address   = "10.0.0.1"
  platform  = 3
  zone_name = "Production"
  node_name = "Network"
  comment   = "Core network switch"

  protocols {
    name = "ssh"
    port = 22
  }
}
```

## Argument Reference

- `name` - (Required) The name of the device.
- `address` - (Required) The address (IP/hostname).
- `platform` - (Required) The platform ID (integer).
- `zone_name` - (Required) The name of the zone.
- `node_name` - (Required) The name of the node.
- `comment` - (Optional) A description.
- `is_active` - (Optional) Whether the device is active. Default: `true`.
- `accounts` - (Optional) List of account blocks.
- `protocols` - (Optional) List of protocol blocks (`name`, `port`).

## Attribute Reference

- `id` - The UUID of the device asset.
- `zone_id` - The resolved zone UUID.
- `node_ids` - The resolved node UUIDs.
