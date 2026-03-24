---
page_title: "jumpserver_account Data Source - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Look up an account by name in JumpServer.
---

# `jumpserver_account` Data Source

Use this data source to look up an existing account by name.

## Example Usage

```hcl
data "jumpserver_account" "root" {
  name = "root"
}
```

## Argument Reference

- **`name`** - (Required) The name of the account to look up.

## Attribute Reference

- **`id`** - The UUID of the account.
- **`username`** - The username of the account.
- **`asset`** - The UUID of the associated asset.
- **`secret_type`** - The type of secret (password, ssh_key, etc.).
- **`is_active`** - Whether the account is active.
