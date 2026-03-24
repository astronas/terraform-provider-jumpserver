---
page_title: "jumpserver_apply_command_ticket Resource - terraform-provider-jumpserver"
subcategory: "Tickets"
description: |-
  Manages an apply-command ticket in JumpServer.
---

# `jumpserver_apply_command_ticket` Resource

The `jumpserver_apply_command_ticket` resource allows you to create tickets for command execution review in JumpServer.

## Example Usage

```hcl
resource "jumpserver_apply_command_ticket" "example" {
  title             = "Run maintenance command"
  apply_run_asset   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  apply_run_account = "root"
  apply_run_command = "systemctl restart nginx"
}
```

## Argument Reference

- **`title`** - (Required) The title of the ticket.
- **`apply_run_asset`** - (Required) The UUID of the asset to run the command on.
- **`apply_run_account`** - (Required) The account to run the command as.
- **`apply_run_command`** - (Required) The command to run.
- **`comment`** - (Optional) A comment for the ticket.

## Attributes Reference

- **`id`** - The UUID of the ticket.
- **`state`** - The ticket state.
- **`status`** - The ticket status.
- **`serial_num`** - The serial number.
- **`org_name`** - The organization name.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
