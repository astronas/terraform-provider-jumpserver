---
page_title: "jumpserver_session Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a session in JumpServer.
---

# `jumpserver_session` Resource

The `jumpserver_session` resource allows you to create and manage sessions in JumpServer.

## Example Usage

```hcl
resource "jumpserver_session" "example" {
  user       = "admin"
  asset      = "web-server-01"
  user_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  asset_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  account    = "root"
  account_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  protocol   = "ssh"
}
```

## Argument Reference

- **`user`** - (Required) The username.
- **`asset`** - (Required) The asset name.
- **`user_id`** - (Required) The user UUID.
- **`asset_id`** - (Required) The asset UUID.
- **`account`** - (Required) The account name.
- **`account_id`** - (Required) The account ID.
- **`protocol`** - (Required) The protocol (e.g. ssh, rdp).
- **`type`** - (Optional) The session type: normal, tunnel, command. Default: `normal`.
- **`login_from`** - (Optional) The login source: ST, RT, WT. Default: `ST`.
- **`remote_addr`** - (Optional) The remote address.
- **`comment`** - (Optional) A comment.

## Attributes Reference

- **`id`** - The UUID of the session.
- **`is_locked`** - Whether the session is locked.
- **`is_finished`** - Whether the session is finished.
- **`is_success`** - Whether the login was successful.
- **`duration`** - The session duration in seconds.
- **`command_amount`** - The number of commands executed.
- **`error_reason`** - The error reason if login failed.
- **`terminal_display`** - The terminal display name.
- **`can_replay`** - Whether the session can be replayed.
- **`can_join`** - Whether the session can be joined.
- **`can_terminate`** - Whether the session can be terminated.
- **`date_start`** - The session start date.
- **`date_end`** - The session end date.
