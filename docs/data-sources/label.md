---
page_title: "jumpserver_label Data Source - terraform-provider-jumpserver"
subcategory: "Labels"
description: |-
  Look up a label by name in JumpServer.
---

# `jumpserver_label` Data Source

Use this data source to look up an existing label by name.

## Example Usage

```hcl
data "jumpserver_label" "env" {
  name = "environment"
}
```

## Argument Reference

- **`name`** - (Required) The name of the label to look up.

## Attribute Reference

- **`id`** - The UUID of the label.
- **`value`** - The value of the label.
- **`color`** - The color of the label.
- **`comment`** - A comment or description.
