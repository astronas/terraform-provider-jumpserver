---
page_title: "jumpserver_system_role Data Source - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Look up a system role in JumpServer by name.
---

# `jumpserver_system_role` Data Source

Use this data source to look up an existing system role in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_system_role" "admin" {
  name = "SystemAdmin"
}
```

## Argument Reference

* `name` - (Required) The name of the system role to look up.

## Attribute Reference

* `id` - The UUID of the system role.
* `comment` - The description of the role.
* `builtin` - Whether this is a built-in role.
