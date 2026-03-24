---
page_title: "jumpserver_labeled_resource Resource - terraform-provider-jumpserver"
subcategory: "Labels"
description: |-
  Manages a labeled resource association in JumpServer.
---

# `jumpserver_labeled_resource` Resource

The `jumpserver_labeled_resource` resource allows you to attach labels to resources in JumpServer.

## Example Usage

```hcl
resource "jumpserver_labeled_resource" "tag" {
  label    = "550e8400-e29b-41d4-a716-446655440000"
  res_type = "asset"
  res_id   = "660e8400-e29b-41d4-a716-446655440001"
}
```

## Argument Reference

* `label` - (Required) The UUID of the label to attach.
* `res_type` - (Required) The type of the resource (e.g., `asset`, `user`).
* `res_id` - (Required) The UUID of the resource to label.

## Attribute Reference

* `id` - The UUID of the labeled resource association.
