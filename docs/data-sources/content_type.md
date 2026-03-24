---
page_title: "jumpserver_content_type Data Source - terraform-provider-jumpserver"
subcategory: "Labels"
description: |-
  Look up a content type (resource type) in JumpServer by name.
---

# `jumpserver_content_type` Data Source

Use this data source to look up an existing content type (resource type) in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_content_type" "example" {
  name = "asset"
}
```

## Argument Reference

- **`name`** - (Required) The name of the content type to look up.

## Attributes Reference

- **`id`** - The ID of the content type.
- **`app_label`** - The application label.
- **`model`** - The model name.
