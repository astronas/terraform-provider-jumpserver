---
page_title: "jumpserver_label Resource"
subcategory: "Assets"
description: |-
  Manages a label (tag) in JumpServer.
---

# jumpserver_label

Manages a label in JumpServer. Labels are key-value tags that can be applied to assets, zones, accounts, and other resources.

## Example Usage

```hcl
resource "jumpserver_label" "env_prod" {
  name    = "environment"
  value   = "production"
  color   = "#00FF00"
  comment = "Production environment tag"
}
```

## Argument Reference

- `name` - (Required) The name (key) of the label.
- `value` - (Required) The value of the label.
- `color` - (Optional) The display color of the label.
- `comment` - (Optional) A description for the label.

## Attribute Reference

- `id` - The UUID of the label.
