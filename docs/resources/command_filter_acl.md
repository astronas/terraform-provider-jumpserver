---
page_title: "jumpserver_command_filter_acl Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages a JumpServer Command Filter ACL.
---

# jumpserver_command_filter_acl

Manages a JumpServer Command Filter ACL. Controls which commands users can execute on assets.

## Example Usage

```hcl
resource "jumpserver_command_group" "dangerous" {
  name    = "dangerous-commands"
  type    = "command"
  content = "rm\nshutdown\nreboot"
}

resource "jumpserver_command_filter_acl" "block_dangerous" {
  name     = "block-dangerous-commands"
  users    = jsonencode({ type = "all" })
  assets   = jsonencode({ type = "all" })
  accounts = ["@ALL"]
  command_groups = [jumpserver_command_group.dangerous.id]
  priority  = 10
  action    = "reject"
  is_active = true
  comment   = "Block dangerous commands for all users"
}
```

## Argument Reference

- `name` - (Required) The name of the command filter ACL.
- `users` - (Required) User filter as JSON string. Example: `jsonencode({ type = "ids", ids = ["uuid1"] })` or `jsonencode({ type = "all" })`.
- `assets` - (Required) Asset filter as JSON string. Same format as users.
- `accounts` - (Required) List of account usernames (e.g. `["root", "@ALL"]`).
- `command_groups` - (Optional) List of command group UUIDs.
- `priority` - (Optional) Priority 1-100. Lower values match first. Defaults to `50`.
- `action` - (Optional) Action when matched: `reject`, `accept`, `review`. Defaults to `reject`.
- `reviewers` - (Optional) List of reviewer user UUIDs. Required when action is `review`.
- `is_active` - (Optional) Whether the ACL is active. Defaults to `true`.
- `comment` - (Optional) A description for the ACL.

## Attribute Reference

- `id` - The UUID of the command filter ACL.

## Import

```shell
terraform import jumpserver_command_filter_acl.example <uuid>
```
