---
page_title: "jumpserver_command_group Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages a JumpServer Command Group.
---

# jumpserver_command_group

Manages a JumpServer Command Group. Command groups define sets of commands that can be used in command filter ACLs.

## Example Usage

```hcl
resource "jumpserver_command_group" "dangerous" {
  name        = "dangerous-commands"
  type        = "command"
  content     = "rm\nshutdown\nreboot\nformat"
  ignore_case = true
  comment     = "Dangerous system commands"
}

resource "jumpserver_command_group" "regex_pattern" {
  name        = "password-commands"
  type        = "regex"
  content     = "passwd.*|chpasswd.*"
  ignore_case = true
}
```

## Argument Reference

- `name` - (Required) The name of the command group.
- `content` - (Required) Commands to match, one per line.
- `type` - (Optional) Type of matching: `command` or `regex`. Defaults to `command`.
- `ignore_case` - (Optional) Whether to ignore case when matching. Defaults to `true`.
- `comment` - (Optional) A description for the command group.

## Attribute Reference

- `id` - The UUID of the command group.

## Import

Command groups can be imported using their UUID:

```shell
terraform import jumpserver_command_group.example <uuid>
```
