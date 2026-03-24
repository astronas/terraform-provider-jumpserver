---
page_title: "jumpserver_node Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages an asset tree node in JumpServer.
---

# `jumpserver_node` Resource

The jumpserver_node resource allows you to create and manage asset nodes (tree structure) in JumpServer. Nodes are used to organize assets hierarchically.

## Example Usage

### Top-level node

```hcl
resource "jumpserver_node" "production" {
  value = "Production"
}
```

### Child node under an existing parent

```hcl
resource "jumpserver_node" "prod_web" {
  value     = "Web Servers"
  parent_id = jumpserver_node.production.id
}
```

## Argument Reference

* `value` - (Required) The display name of the node.
* `parent_id` - (Optional) The ID of the parent node. If omitted, the node is created under the default root.

## Attribute Reference

* `id` - The ID of the node.
* `value` - The display name of the node.
* `full_value` - The full path of the node in the tree (e.g. `/Default/Production/Web Servers`).
* `key` - The tree key of the node (e.g. `1:1:2`).
* `child_mark` - The child mark counter for the node.
