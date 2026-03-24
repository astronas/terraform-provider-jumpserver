---
page_title: "jumpserver_activity_log Data Source - terraform-provider-jumpserver"
subcategory: "Audits"
description: |-
  Look up activity logs in JumpServer by resource ID.
---

# `jumpserver_activity_log` Data Source

Use this data source to look up activity log entries for a specific resource in JumpServer.

## Example Usage

```hcl
data "jumpserver_activity_log" "recent" {
  resource_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

* `resource_id` - (Required) The resource ID to look up activity logs for.

## Attribute Reference

* `id` - The ID of the activity log entry.
* `timestamp` - The activity timestamp.
* `detail_url` - The detail URL for the activity.
* `content` - The activity content.
* `r_type` - The resource type.
