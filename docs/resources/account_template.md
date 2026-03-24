---
page_title: "jumpserver_account_template Resource - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Manages a JumpServer Account Template.
---

# jumpserver_account_template

Manages a JumpServer Account Template. Templates define reusable account configurations that can be applied across multiple assets.

## Example Usage

```hcl
resource "jumpserver_account_template" "root_password" {
  name            = "root-password-template"
  username        = "root"
  secret_type     = "password"
  secret_strategy = "random"
  privileged      = true
  auto_push       = true

  password_rules {
    length    = 24
    lowercase = true
    uppercase = true
    digit     = true
    symbol    = false
  }

  comment = "Root account with random password"
}

resource "jumpserver_account_template" "deploy_key" {
  name            = "deploy-ssh-key"
  username        = "deploy"
  secret_type     = "ssh_key"
  secret_strategy = "specific"
  secret          = file("~/.ssh/deploy_key")
  passphrase      = var.deploy_key_passphrase
  auto_push       = true
}
```

## Argument Reference

- `name` - (Required) The name of the account template.
- `username` - (Optional) The username. For AD accounts use `username@domain`.
- `secret_type` - (Optional) Type of secret: `password` or `ssh_key`. Defaults to `password`.
- `secret` - (Optional, Sensitive) The secret value (password or SSH key content).
- `passphrase` - (Optional, Sensitive) Passphrase for SSH key.
- `secret_strategy` - (Optional) Secret strategy: `specific` (user-provided) or `random` (auto-generated). Defaults to `specific`.
- `password_rules` - (Optional) Password complexity rules block (used with `random` strategy):
  - `length` - (Optional) Password length, 8-36. Defaults to `16`.
  - `lowercase` - (Optional) Include lowercase letters. Defaults to `true`.
  - `uppercase` - (Optional) Include uppercase letters. Defaults to `true`.
  - `digit` - (Optional) Include digits. Defaults to `true`.
  - `symbol` - (Optional) Include symbols. Defaults to `true`.
  - `exclude_symbols` - (Optional) Symbols to exclude from generated passwords.
- `platforms` - (Optional) List of platform IDs to associate.
- `su_from` - (Optional) Account template ID to switch from (su).
- `privileged` - (Optional) Whether this is a privileged account. Defaults to `false`.
- `is_active` - (Optional) Whether the template is active. Defaults to `true`.
- `auto_push` - (Optional) Automatically push account to assets. Defaults to `false`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the account template.

## Import

```shell
terraform import jumpserver_account_template.example <uuid>
```
