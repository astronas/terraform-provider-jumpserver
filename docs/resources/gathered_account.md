---
page_title: "jumpserver_gathered_account Resource - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Manages a gathered account in JumpServer.
---

# `jumpserver_gathered_account` Resource

The `jumpserver_gathered_account` resource allows you to manage gathered (discovered) accounts on assets in JumpServer.

## Example Usage

```hcl
resource "jumpserver_gathered_account" "example" {
  asset = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- **`asset`** - (Required) The UUID of the asset this gathered account belongs to.

## Attributes Reference

- **`id`** - The UUID of the gathered account.
- **`username`** - The discovered username.
- **`address_last_login`** - The last login address.
- **`status`** - The account status.
- **`remote_present`** - Whether the account is present remotely.
- **`present`** - Whether the account is present.
- **`date_last_login`** - The date of last login.
- **`date_updated`** - The last update date.
