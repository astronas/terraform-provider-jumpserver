---
page_title: "jumpserver_apply_login_asset_ticket Resource - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Manages an apply-login-asset ticket in JumpServer.
---

# `jumpserver_apply_login_asset_ticket` Resource

The `jumpserver_apply_login_asset_ticket` resource allows you to create tickets for login-to-asset review in JumpServer.

## Example Usage

```hcl
resource "jumpserver_apply_login_asset_ticket" "example" {
  title               = "Login asset review"
  apply_login_user    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  apply_login_asset   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  apply_login_account = "root"
}
```

## Argument Reference

- **`title`** - (Required) The title of the ticket.
- **`apply_login_user`** - (Optional) The UUID of the user logging in.
- **`apply_login_asset`** - (Optional) The UUID of the asset being logged into.
- **`apply_login_account`** - (Optional) The account name used for login.
- **`comment`** - (Optional) A comment for the ticket.

## Attributes Reference

- **`id`** - The UUID of the ticket.
- **`state`** - The ticket state.
- **`status`** - The ticket status.
- **`serial_num`** - The serial number.
- **`org_name`** - The organization name.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
