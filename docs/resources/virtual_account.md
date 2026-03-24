---
page_title: "jumpserver_virtual_account Resource - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Manages a virtual account in JumpServer.
---

# `jumpserver_virtual_account` Resource

The `jumpserver_virtual_account` resource allows you to create and manage virtual accounts used for special account selection during asset connections.

## Example Usage

```hcl
resource "jumpserver_virtual_account" "input" {
  alias              = "@INPUT"
  secret_from_login  = false
}
```

## Argument Reference

* `alias` - (Required) The alias for the virtual account. Valid values: `@INPUT`, `@USER`, `@ANON`, `@SPEC`.
* `secret_from_login` - (Optional) Whether to use the login secret. Default: `false`.

## Attribute Reference

* `id` - The UUID of the virtual account.
* `username` - The resolved username.
* `name` - The display name.
* `comment` - The description.
