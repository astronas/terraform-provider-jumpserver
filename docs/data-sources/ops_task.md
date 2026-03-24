---
page_title: "jumpserver_ops_task Data Source - terraform-provider-jumpserver"
subcategory: "Ops"
description: |-
  Look up an ops task in JumpServer by name.
---

# `jumpserver_ops_task` Data Source

Use this data source to look up a Celery ops task in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_ops_task" "backup" {
  name = "backup_task"
}
```

## Argument Reference

* `name` - (Required) The name of the ops task to look up.

## Attribute Reference

* `id` - The UUID of the ops task.
* `meta` - Task metadata as a map of strings.
