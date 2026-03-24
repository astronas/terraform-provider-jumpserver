---
page_title: "jumpserver_account_backup_plan Resource - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Manages a JumpServer Account Backup Plan.
---

# jumpserver_account_backup_plan

Manages a JumpServer Account Backup Plan. Automates periodic backup of account secrets.

## Example Usage

```hcl
resource "jumpserver_account_backup_plan" "daily" {
  name        = "daily-backup"
  is_periodic = true
  interval    = 24
  is_active   = true

  recipients_part_one = [jumpserver_user.admin.id]

  comment = "Daily backup of all account secrets"
}

resource "jumpserver_account_backup_plan" "cron_backup" {
  name        = "weekly-backup"
  is_periodic = true
  crontab     = "0 2 * * 0"
  is_active   = true

  zip_encrypt_password = var.backup_password

  recipients_part_one = [jumpserver_user.admin.id]
  recipients_part_two = [jumpserver_user.security.id]

  comment = "Weekly encrypted backup"
}
```

## Argument Reference

- `name` - (Required) The name of the backup plan.
- `is_periodic` - (Optional) Whether to run periodically. Defaults to `true`.
- `interval` - (Optional) Execution interval in hours (1-65535). Defaults to `24`.
- `crontab` - (Optional) Cron expression for scheduling (alternative to interval).
- `accounts` - (Optional) List of account IDs to back up. Empty means all.
- `types` - (Optional) Asset type filters.
- `nodes` - (Optional) List of node IDs to include.
- `recipients_part_one` - (Optional) User IDs for email recipients (part one).
- `recipients_part_two` - (Optional) User IDs for email recipients (part two).
- `obj_recipients_part_one` - (Optional) User IDs for object storage recipients (part one).
- `obj_recipients_part_two` - (Optional) User IDs for object storage recipients (part two).
- `zip_encrypt_password` - (Optional, Sensitive) Password to encrypt the backup zip file.
- `is_password_divided_by_email` - (Optional) Send password separately by email. Defaults to `true`.
- `is_password_divided_by_obj_storage` - (Optional) Send password separately to object storage. Defaults to `true`.
- `is_active` - (Optional) Whether the plan is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the backup plan.

## Import

```shell
terraform import jumpserver_account_backup_plan.example <uuid>
```
