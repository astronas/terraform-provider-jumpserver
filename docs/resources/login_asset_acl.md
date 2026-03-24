---
page_title: "jumpserver_login_asset_acl Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages a JumpServer Login Asset ACL.
---

# jumpserver_login_asset_acl

Manages a JumpServer Login Asset ACL. Controls user access to assets based on rules (IP, time, etc.).

## Example Usage

```hcl
resource "jumpserver_login_asset_acl" "prod_restrict" {
  name     = "restrict-prod-access"
  users    = jsonencode({ type = "all" })
  assets   = jsonencode({ type = "ids", ids = [jumpserver_host.prod.id] })
  accounts = ["root"]
  rules    = jsonencode({ ip_group = ["10.0.0.0/8"] })
  priority = 10
  action   = "review"
  reviewers = [jumpserver_user.admin.id]
  is_active = true
  comment  = "Require review for prod root access"
}
```

## Argument Reference

- `name` - (Required) The name of the login asset ACL.
- `users` - (Required) User filter as JSON string.
- `assets` - (Required) Asset filter as JSON string.
- `accounts` - (Required) List of account usernames.
- `rules` - (Required) Rules as JSON string (flexible structure for conditions).
- `priority` - (Optional) Priority 1-100. Defaults to `50`.
- `action` - (Optional) Action: `reject`, `accept`, `review`, `notice`. Defaults to `reject`.
- `reviewers` - (Optional) List of reviewer user UUIDs.
- `is_active` - (Optional) Whether the ACL is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the login asset ACL.

## Import

```shell
terraform import jumpserver_login_asset_acl.example <uuid>
```
