---
page_title: "jumpserver_host Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a host by name in JumpServer.
---

# `jumpserver_host` Data Source

Use this data source to look up an existing host by name.

## Example Usage

```hcl
data "jumpserver_host" "example" {
  name = "server-lxc1"
}
```

## Argument Reference

- **`name`** - (Required) The name of the host to look up.

## Attribute Reference

- **`id`** - The UUID of the host.
- **`address`** - The address (IP or hostname) of the host.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the host is active.
