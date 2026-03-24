---
page_title: "jumpserver_protocol Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a protocol in JumpServer by name.
---

# `jumpserver_protocol` Data Source

Use this data source to look up an existing protocol in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_protocol" "ssh" {
  name = "ssh"
}
```

## Argument Reference

- **`name`** - (Required) The name of the protocol to look up.

## Attributes Reference

- **`id`** - The ID of the protocol.
- **`port`** - The default port of the protocol.
