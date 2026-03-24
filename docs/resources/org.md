---
page_title: "jumpserver_org Resource - terraform-provider-jumpserver"
subcategory: "RBAC"
description: |-
  Manages an organization in JumpServer.
---

# `jumpserver_org` Resource

The `jumpserver_org` resource allows you to create and manage organizations in JumpServer for multi-tenancy.

## Example Usage

```hcl
resource "jumpserver_org" "dev" {
  name    = "Development"
  comment = "Development organization"
}
```

## Argument Reference

* `name` - (Required) The name of the organization.
* `comment` - (Optional) A description for the organization.

## Attribute Reference

* `id` - The UUID of the organization.
* `is_default` - Whether this is the default organization.
* `is_root` - Whether this is the root organization.
