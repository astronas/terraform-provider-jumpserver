---
page_title: "jumpserver_org Data Source - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Look up an organization in JumpServer by name.
---

# `jumpserver_org` Data Source

Use this data source to look up an existing organization in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_org" "default" {
  name = "Default"
}
```

## Argument Reference

* `name` - (Required) The name of the organization to look up.

## Attribute Reference

* `id` - The UUID of the organization.
* `comment` - The description of the organization.
