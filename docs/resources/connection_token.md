---
page_title: "jumpserver_connection_token Resource - terraform-provider-jumpserver"
subcategory: "Authentication"
description: |-
  Manages a connection token in JumpServer.
---

# `jumpserver_connection_token` Resource

The `jumpserver_connection_token` resource allows you to create and manage connection tokens for asset access in JumpServer.

## Example Usage

```hcl
resource "jumpserver_connection_token" "example" {
  account        = "root"
  connect_method = "web_cli"
  protocol       = "ssh"
  asset          = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`account`** - (Required) The account name.
- **`connect_method`** - (Required) The connection method.
- **`protocol`** - (Optional) The protocol (e.g. ssh, rdp). Default: `ssh`.
- **`asset`** - (Optional) The asset UUID.
- **`input_username`** - (Optional) The input username.
- **`input_secret`** - (Optional, Sensitive) The input secret (write-only).
- **`remote_addr`** - (Optional) The remote address.
- **`is_active`** - (Optional) Whether the token is active. Default: `true`.
- **`is_reusable`** - (Optional) Whether the token is reusable. Default: `false`.

## Attributes Reference

- **`id`** - The UUID of the connection token.
- **`value`** - (Sensitive) The token value.
- **`user_display`** - The user display name.
- **`asset_display`** - The asset display name.
- **`is_expired`** - Whether the token is expired.
- **`expire_time`** - The expiration time in seconds.
- **`date_expired`** - The expiration date.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
