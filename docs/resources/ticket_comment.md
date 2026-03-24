---
page_title: "jumpserver_ticket_comment Resource - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Manages a ticket comment in JumpServer.
---

# `jumpserver_ticket_comment` Resource

The `jumpserver_ticket_comment` resource allows you to create and manage comments on tickets in JumpServer.

## Example Usage

```hcl
resource "jumpserver_ticket_comment" "example" {
  body = "This ticket has been reviewed and approved."
}
```

## Argument Reference

- **`body`** - (Required) The body text of the comment.

## Attributes Reference

- **`id`** - The UUID of the comment.
- **`user_display`** - The display name of the commenter.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
