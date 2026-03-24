---
page_title: "jumpserver_login_acl Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages a JumpServer Login ACL.
---

# jumpserver_login_acl

Manages a JumpServer Login ACL. Controls user login access based on IP, time, and MFA conditions.

## Example Usage

```hcl
resource "jumpserver_login_acl" "restrict_ip" {
  name     = "restrict-login-ip"
  users    = jsonencode({ type = "all" })
  rules    = jsonencode({ ip_group = ["192.168.1.0/24"] })
  priority = 10
  action   = "accept"
  is_active = true
  comment  = "Allow login only from internal network"
}
```

## Argument Reference

- `name` - (Required) The name of the login ACL.
- `users` - (Required) User filter as JSON string. Example: `jsonencode({ type = "all" })`.
- `rules` - (Required) Rules as JSON string (flexible structure for IP, time, MFA conditions).
- `priority` - (Optional) Priority 1-100. Defaults to `50`.
- `action` - (Optional) Action: `reject`, `accept`, `review`, `notice`. Defaults to `reject`.
- `reviewers` - (Optional) List of reviewer user UUIDs.
- `is_active` - (Optional) Whether the ACL is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the login ACL.

## Import

```shell
terraform import jumpserver_login_acl.example <uuid>
```
