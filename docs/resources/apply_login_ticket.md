---
page_title: "jumpserver_apply_login_ticket Resource - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Manages an apply-login ticket in JumpServer.
---

# `jumpserver_apply_login_ticket` Resource

The `jumpserver_apply_login_ticket` resource allows you to create tickets for login review in JumpServer.

## Example Usage

```hcl
resource "jumpserver_apply_login_ticket" "example" {
  title               = "Login review request"
  apply_login_ip      = "192.168.1.100"
  apply_login_city    = "Paris"
  apply_login_datetime = "2024-01-15T10:30:00Z"
}
```

## Argument Reference

- **`title`** - (Required) The title of the ticket.
- **`apply_login_ip`** - (Optional) The login IP address.
- **`apply_login_city`** - (Optional) The login city.
- **`apply_login_datetime`** - (Optional) The login date/time.
- **`comment`** - (Optional) A comment for the ticket.

## Attributes Reference

- **`id`** - The UUID of the ticket.
- **`state`** - The ticket state.
- **`status`** - The ticket status.
- **`serial_num`** - The serial number.
- **`org_name`** - The organization name.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
