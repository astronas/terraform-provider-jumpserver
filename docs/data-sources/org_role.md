---
page_title: "jumpserver_org_role Data Source - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Look up an organization role in JumpServer by name.
---

# `jumpserver_org_role` Data Source

Use this data source to look up an existing organization role in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_org_role" "auditor" {
  name = "OrgAuditor"
}
```

## Argument Reference

* `name` - (Required) The name of the organization role to look up.

## Attribute Reference

* `id` - The UUID of the organization role.
* `comment` - The description of the role.
* `builtin` - Whether this is a built-in role.
