---
page_title: "jumpserver_platform Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a platform in JumpServer by name.
---

# `jumpserver_platform` Data Source

Use this data source to look up an existing platform in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_platform" "linux" {
  name = "Linux"
}

resource "jumpserver_host" "example" {
  name     = "my-server"
  address  = "10.0.0.1"
  platform = data.jumpserver_platform.linux.id
  # ...
}
```

## Argument Reference

* `name` - (Required) The name of the platform to look up.

## Attribute Reference

* `id` - The ID of the platform (numeric, returned as string).
* `comment` - The description of the platform.
* `category` - The category of the platform (e.g., host, device).
* `type` - The type of the platform.
