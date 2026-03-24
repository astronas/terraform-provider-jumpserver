---
page_title: "jumpserver_access_key Resource - terraform-provider-jumpserver"
subcategory: "Authentication"
description: |-
  Manages an access key in JumpServer.
---

# `jumpserver_access_key` Resource

The `jumpserver_access_key` resource allows you to create and manage API access keys for authentication.

## Example Usage

```hcl
resource "jumpserver_access_key" "api_key" {
  is_active = true
  ip_group  = ["10.0.0.0/8", "172.16.0.0/12"]
}
```

## Argument Reference

* `is_active` - (Optional) Whether the access key is active. Default: `true`.
* `ip_group` - (Optional) List of allowed IP addresses or CIDR ranges. Default: `["*"]`.

## Attribute Reference

* `id` - The UUID of the access key.
* `date_created` - The creation date.
* `date_last_used` - The date of last use.
