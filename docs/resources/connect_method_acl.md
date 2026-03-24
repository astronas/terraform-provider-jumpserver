---
page_title: "jumpserver_connect_method_acl Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages a JumpServer Connect Method ACL.
---

# jumpserver_connect_method_acl

Manages a JumpServer Connect Method ACL. Controls which connection methods users can use (SSH, RDP, etc.).

## Example Usage

```hcl
resource "jumpserver_connect_method_acl" "ssh_only" {
  name            = "allow-ssh-only"
  users           = jsonencode({ type = "all" })
  connect_methods = ["ssh"]
  priority        = 10
  action          = "accept"
  is_active       = true
  comment         = "Only allow SSH connections"
}
```

## Argument Reference

- `name` - (Required) The name of the connect method ACL.
- `users` - (Required) User filter as JSON string.
- `connect_methods` - (Optional) List of allowed connect method identifiers.
- `priority` - (Optional) Priority 1-100. Defaults to `50`.
- `action` - (Optional) Action: `reject`, `accept`, `review`. Defaults to `reject`.
- `reviewers` - (Optional) List of reviewer user UUIDs.
- `is_active` - (Optional) Whether the ACL is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the connect method ACL.

## Import

```shell
terraform import jumpserver_connect_method_acl.example <uuid>
```
