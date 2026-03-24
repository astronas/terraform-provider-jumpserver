---
page_title: "jumpserver_account_template Data Source - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Look up an account template by name in JumpServer.
---

# `jumpserver_account_template` Data Source

Use this data source to look up an existing account template by name.

## Example Usage

```hcl
data "jumpserver_account_template" "ssh_template" {
  name = "ssh-root-template"
}
```

## Argument Reference

- **`name`** - (Required) The name of the account template to look up.

## Attribute Reference

- **`id`** - The UUID of the account template.
- **`username`** - The username defined in the template.
- **`secret_type`** - The type of secret (password, ssh_key, etc.).
- **`is_active`** - Whether the account template is active.
- **`comment`** - A comment or description.
