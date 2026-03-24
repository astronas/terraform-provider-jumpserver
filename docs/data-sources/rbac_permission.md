---
page_title: "jumpserver_rbac_permission Data Source - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Look up an RBAC permission in JumpServer by name.
---

# `jumpserver_rbac_permission` Data Source

Use this data source to look up an RBAC permission in JumpServer by its name. Useful for building permission lists for role assignments.

## Example Usage

```hcl
data "jumpserver_rbac_permission" "view_asset" {
  name = "Can view asset"
}
```

## Argument Reference

* `name` - (Required) The name of the permission to look up.

## Attribute Reference

* `id` - The integer ID of the permission.
* `codename` - The codename of the permission.
* `content_type` - The content type ID.
