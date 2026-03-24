---
page_title: "jumpserver_ticket_flow Resource - terraform-provider-jumpserver"
subcategory: "Ops"
description: |-
  Manages a ticket approval flow in JumpServer.
---

# `jumpserver_ticket_flow` Resource

The `jumpserver_ticket_flow` resource allows you to create and manage ticket approval flows (workflows) in JumpServer.

## Example Usage

```hcl
resource "jumpserver_ticket_flow" "login_approval" {
  type           = "login_confirm"
  approval_level = 1

  rules = jsonencode([
    {
      assignees = ["550e8400-e29b-41d4-a716-446655440000"]
      type      = "super_admin"
    }
  ])
}
```

## Argument Reference

* `type` - (Required, ForceNew) The type of ticket flow. Valid values: `apply_asset`, `login_confirm`, `login_asset_confirm`, `command_confirm`.
* `approval_level` - (Required) The approval level (1-3).
* `rules` - (Required) JSON string of approval rules defining assignees and conditions.

## Attribute Reference

* `id` - The UUID of the ticket flow.
