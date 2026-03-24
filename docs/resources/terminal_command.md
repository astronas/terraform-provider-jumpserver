---
page_title: "jumpserver_terminal_command Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a terminal command record in JumpServer.
---

# `jumpserver_terminal_command` Resource

The `jumpserver_terminal_command` resource allows you to manage terminal command records in JumpServer.

## Example Usage

```hcl
resource "jumpserver_terminal_command" "example" {
  user      = "admin"
  asset     = "web-server-01"
  input     = "ls -la"
  session   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  account   = "root"
  output    = "total 64..."
  timestamp = 1704067200
}
```

## Argument Reference

- **`user`** - (Required) The username.
- **`asset`** - (Required) The asset name.
- **`input`** - (Required) The command input.
- **`session`** - (Required) The session UUID.
- **`account`** - (Required) The account name.
- **`output`** - (Required) The command output.
- **`timestamp`** - (Required) The timestamp of the command.
- **`risk_level`** - (Optional) The risk level: 0 (Ordinary), 1 (Dangerous). Default: `0`.
- **`org_id`** - (Optional) The organization ID.

## Attributes Reference

- **`id`** - The UUID of the command record.
- **`timestamp_display`** - The formatted timestamp display.
- **`remote_addr`** - The remote address.
