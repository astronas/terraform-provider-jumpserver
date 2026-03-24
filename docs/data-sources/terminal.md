---
page_title: "jumpserver_terminal Data Source - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Look up a terminal in JumpServer by name.
---

# `jumpserver_terminal` Data Source

Use this data source to look up an existing terminal in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_terminal" "example" {
  name = "koko-01"
}
```

## Argument Reference

- **`name`** - (Required) The name of the terminal to look up.

## Attributes Reference

- **`id`** - The UUID of the terminal.
- **`remote_addr`** - The remote address.
- **`type`** - The terminal type.
- **`is_active`** - Whether the terminal is active.
- **`is_alive`** - Whether the terminal is alive.
- **`comment`** - A comment or description.
