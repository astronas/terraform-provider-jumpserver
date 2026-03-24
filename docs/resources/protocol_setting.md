---
page_title: "jumpserver_protocol_setting Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a protocol setting in JumpServer.
---

# `jumpserver_protocol_setting` Resource

Manages protocol settings for asset connections in JumpServer.

## Example Usage

```hcl
resource "jumpserver_protocol_setting" "ssh" {
  name    = "ssh"
  port    = 22
  primary = true
  public  = true
}
```

## Argument Reference

- **`name`** - (Required) The name of the protocol.
- **`port`** - (Optional) The default port for the protocol. Defaults to `0`.
- **`primary`** - (Optional) Whether this is a primary protocol.
- **`required`** - (Optional) Whether this protocol is required.
- **`default`** - (Optional) Whether this is a default protocol.
- **`public`** - (Optional) Whether this protocol is public.
- **`setting`** - (Optional) Additional settings as a JSON string.

## Attribute Reference

- **`id`** - The ID of the protocol setting.
- **`port_from_addr`** - Whether the port is derived from the address.
