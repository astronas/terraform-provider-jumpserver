---
page_title: "jumpserver_apply_asset_ticket Resource - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Manages an apply-asset ticket in JumpServer.
---

# `jumpserver_apply_asset_ticket` Resource

The `jumpserver_apply_asset_ticket` resource allows you to create tickets for requesting asset access in JumpServer.

## Example Usage

```hcl
resource "jumpserver_apply_asset_ticket" "example" {
  title              = "Request access to production"
  apply_date_start   = "2024-01-01T00:00:00Z"
  apply_date_expired = "2024-12-31T23:59:59Z"
  comment            = "Need production access for deployment"
}
```

## Argument Reference

- **`title`** - (Required) The title of the ticket.
- **`apply_nodes`** - (Optional) List of node UUIDs to request access to.
- **`apply_assets`** - (Optional) List of asset UUIDs to request access to.
- **`apply_accounts`** - (Optional) List of account names to request.
- **`apply_actions`** - (Optional) List of actions to request.
- **`apply_date_start`** - (Optional) The start date for the requested access.
- **`apply_date_expired`** - (Optional) The expiration date for the requested access.
- **`comment`** - (Optional) A comment for the ticket.

## Attributes Reference

- **`id`** - The UUID of the ticket.
- **`state`** - The ticket state.
- **`status`** - The ticket status.
- **`serial_num`** - The serial number.
- **`apply_permission_name`** - The applied permission name.
- **`org_name`** - The organization name.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
