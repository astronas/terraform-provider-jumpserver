---
page_title: "jumpserver_terminal Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a terminal in JumpServer.
---

# `jumpserver_terminal` Resource

The `jumpserver_terminal` resource allows you to create and manage terminals (component instances) in JumpServer.

## Example Usage

```hcl
resource "jumpserver_terminal" "example" {
  name            = "koko-01"
  remote_addr     = "10.0.0.10"
  command_storage = "default"
  replay_storage  = "default"
  comment         = "Primary terminal"
}
```

## Argument Reference

- **`name`** - (Required) The name of the terminal.
- **`remote_addr`** - (Optional) The remote address.
- **`command_storage`** - (Optional) The command storage backend name.
- **`replay_storage`** - (Optional) The replay storage backend name.
- **`comment`** - (Optional) A comment or description.

## Attributes Reference

- **`id`** - The UUID of the terminal.
- **`type`** - The terminal type.
- **`is_active`** - Whether the terminal is active.
- **`is_alive`** - Whether the terminal is alive.
- **`session_online`** - The number of online sessions.
- **`date_created`** - The creation date.
