---
page_title: "jumpserver_operate_log Data Source - terraform-provider-jumpserver"
subcategory: "Audits"
description: |-
  Look up operate logs in JumpServer by user.
---

# `jumpserver_operate_log` Data Source

Use this data source to look up the most recent operation log entry for a user in JumpServer.

## Example Usage

```hcl
data "jumpserver_operate_log" "admin_ops" {
  user = "admin"
}
```

## Argument Reference

* `user` - (Required) The username to look up operate logs for.

## Attribute Reference

* `id` - The UUID of the operate log entry.
* `action` - The action performed.
* `resource_type` - The type of resource affected.
* `resource` - The resource affected.
* `remote_addr` - The remote address.
* `org_id` - The organization ID.
* `org_name` - The organization name.
* `datetime` - The date and time of the operation.
