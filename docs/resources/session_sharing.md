---
page_title: "jumpserver_session_sharing Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a session sharing in JumpServer.
---

# `jumpserver_session_sharing` Resource

The `jumpserver_session_sharing` resource allows you to share sessions with other users in JumpServer.

## Example Usage

```hcl
resource "jumpserver_session_sharing" "example" {
  session      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  expired_time = 60
  is_active    = true
}
```

## Argument Reference

- **`session`** - (Required) The session UUID to share.
- **`expired_time`** - (Optional) Expiration time in minutes (0 = never). Default: `0`.
- **`origin`** - (Optional) The origin URI.
- **`is_active`** - (Optional) Whether the sharing is active. Default: `true`.

## Attributes Reference

- **`id`** - The UUID of the session sharing.
- **`verify_code`** - The verification code.
- **`url`** - The sharing URL.
- **`users_display`** - Display names of users.
- **`action_permission`** - The action permission level.
- **`created_by`** - The creator.
- **`org_id`** - The organization ID.
- **`org_name`** - The organization name.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
