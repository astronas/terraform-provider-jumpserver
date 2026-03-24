---
page_title: "jumpserver_ticket Resource - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Manages a ticket in JumpServer.
---

# `jumpserver_ticket` Resource

The `jumpserver_ticket` resource allows you to create and manage tickets in JumpServer.

## Example Usage

```hcl
resource "jumpserver_ticket" "example" {
  title   = "Access request"
  comment = "Need access to production servers"
}
```

## Argument Reference

- **`title`** - (Required) The title of the ticket.
- **`comment`** - (Optional) A comment for the ticket.

## Attributes Reference

- **`id`** - The UUID of the ticket.
- **`type`** - The ticket type.
- **`state`** - The ticket state.
- **`status`** - The ticket status.
- **`serial_num`** - The serial number of the ticket.
- **`org_id`** - The organization ID.
- **`org_name`** - The organization name.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
