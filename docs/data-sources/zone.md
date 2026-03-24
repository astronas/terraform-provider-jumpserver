---
page_title: "jumpserver_zone Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a zone in JumpServer by name.
---

# `jumpserver_zone` Data Source

Use this data source to look up an existing zone in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_zone" "production" {
  name = "Production"
}
```

## Argument Reference

* `name` - (Required) The name of the zone to look up.

## Attribute Reference

* `id` - The UUID of the zone.
* `comment` - The description of the zone.
