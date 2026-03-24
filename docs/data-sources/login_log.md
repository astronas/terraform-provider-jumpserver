---
page_title: "jumpserver_login_log Data Source - terraform-provider-jumpserver"
subcategory: "Audits"
description: |-
  Look up login logs in JumpServer by username.
---

# `jumpserver_login_log` Data Source

Use this data source to look up the most recent login log entry for a user in JumpServer.

## Example Usage

```hcl
data "jumpserver_login_log" "admin" {
  username = "admin"
}
```

## Argument Reference

* `username` - (Required) The username to look up login logs for.

## Attribute Reference

* `id` - The UUID of the login log entry.
* `type` - The login type.
* `ip` - The IP address.
* `city` - The city of origin.
* `user_agent` - The user agent string.
* `mfa` - The MFA status.
* `reason` - The login reason.
* `reason_display` - The human-readable login reason.
* `backend` - The authentication backend.
* `backend_display` - The human-readable backend name.
* `status` - The login status.
* `datetime` - The login date and time.
