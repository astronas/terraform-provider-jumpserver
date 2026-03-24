---
page_title: "jumpserver_ticket_flow Data Source - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Look up a ticket flow in JumpServer by name.
---

# `jumpserver_ticket_flow` Data Source

Use this data source to look up an existing ticket flow in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_ticket_flow" "example" {
  name = "default-approval-flow"
}
```

## Argument Reference

- **`name`** - (Required) The name of the ticket flow to look up.

## Attributes Reference

- **`id`** - The UUID of the ticket flow.
- **`type`** - The ticket flow type.
- **`is_active`** - Whether the ticket flow is active.
- **`comment`** - A comment or description.
