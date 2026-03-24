---
page_title: "jumpserver_account Resource"
subcategory: "Accounts"
description: |-
  Manages an account in JumpServer.
---

# jumpserver_account

Manages an account in JumpServer. Accounts represent credentials used to access assets (replaces the deprecated system_user resource).

## Example Usage

```hcl
resource "jumpserver_account" "ssh_root" {
  name        = "root"
  username    = "root"
  secret_type = "password"
  secret      = "my-secure-password"
  asset       = jumpserver_host.web_server.id
  privileged  = true
  is_active   = true
  push_now    = true
  comment     = "Root account for web server"
}

resource "jumpserver_account" "ssh_key_deploy" {
  name        = "deploy"
  username    = "deploy"
  secret_type = "ssh_key"
  secret      = file("~/.ssh/id_rsa")
  passphrase  = "my-key-passphrase"
  asset       = jumpserver_host.web_server.id
  privileged  = false
}
```

## Argument Reference

- `name` - (Required) The name of the account.
- `username` - (Required) The username for authentication. Enter null if no username is required. For AD accounts use `username@domain`.
- `secret_type` - (Optional) The type of secret. Valid values: `password`, `ssh_key`, `access_key`, `token`, `api_key`. Default: `password`.
- `secret` - (Optional, Sensitive) The secret value (password, SSH key content, etc.).
- `passphrase` - (Optional, Sensitive) Passphrase for SSH key (only used when `secret_type` is `ssh_key`).
- `asset` - (Required, ForceNew) The ID of the asset this account belongs to. Changing this forces a new resource.
- `template` - (Optional, ForceNew) Account template ID used to create the account from an existing template.
- `privileged` - (Optional) Whether this is a privileged account. Default: `false`.
- `is_active` - (Optional) Whether the account is active. Default: `true`.
- `su_from` - (Optional) The ID of the account to switch from (su).
- `on_invalid` - (Optional) Policy when account already exists: `error`, `skip`, `update`. Default: `error`.
- `push_now` - (Optional) Whether to push the account to the asset immediately after creation. Default: `false`.
- `comment` - (Optional) A description for the account.

## Attribute Reference

- `id` - The UUID of the account.
