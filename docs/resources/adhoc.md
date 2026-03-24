---
page_title: "jumpserver_adhoc Resource - terraform-provider-jumpserver"
subcategory: "Ops"
description: |-
  Manages an ad-hoc command in JumpServer.
---

# `jumpserver_adhoc` Resource

The `jumpserver_adhoc` resource allows you to create and manage ad-hoc commands that can be executed on assets.

## Example Usage

```hcl
resource "jumpserver_adhoc" "check_disk" {
  name    = "Check Disk Usage"
  module  = "shell"
  args    = "df -h"
  comment = "Check disk space usage"
}
```

## Argument Reference

* `name` - (Required) The name of the ad-hoc command.
* `module` - (Required) The module to use. Valid values: `shell`, `winshell`, `python`, `raw`.
* `args` - (Required) The command arguments to execute.
* `comment` - (Optional) A description for the ad-hoc command.

## Attribute Reference

* `id` - The UUID of the ad-hoc command.
* `org_id` - The organization ID.
* `org_name` - The organization name.
* `created_by` - The creator of the ad-hoc command.
