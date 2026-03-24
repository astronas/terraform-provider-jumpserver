---
page_title: "jumpserver_node Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a node in JumpServer by value.
---

# `jumpserver_node` Data Source

Use this data source to look up an existing node in JumpServer by its value (path).

## Example Usage

```hcl
data "jumpserver_node" "linux" {
  value = "Linux Servers"
}
```

## Argument Reference

* `value` - (Required) The value (path) of the node to look up.

## Attribute Reference

* `id` - The UUID of the node.
* `name` - The display name of the node.
* `full_value` - The full path value of the node.
