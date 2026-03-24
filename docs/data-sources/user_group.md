---
page_title: "jumpserver_user_group Data Source - terraform-provider-jumpserver"
subcategory: "Users & Groups"
description: |-
  Look up a user group in JumpServer by name.
---

# `jumpserver_user_group` Data Source

Use this data source to look up an existing user group in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_user_group" "devops" {
  name = "DevOps"
}
```

## Argument Reference

* `name` - (Required) The name of the user group to look up.

## Attribute Reference

* `id` - The UUID of the user group.
* `comment` - The description of the user group.
