---
page_title: "jumpserver_platform Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a platform in JumpServer.
---

# `jumpserver_platform` Resource

The `jumpserver_platform` resource allows you to create and manage custom platforms in JumpServer. Platforms define asset categories, types, and protocol configurations.

## Example Usage

```hcl
resource "jumpserver_platform" "custom_linux" {
  name     = "Custom Linux"
  category = "host"
  type     = "linux"
  charset  = "utf-8"
  comment  = "Custom Linux platform"

  protocols = jsonencode([
    {
      name = "ssh"
      port = 22
    }
  ])
}
```

## Argument Reference

* `name` - (Required) The name of the platform.
* `category` - (Required, ForceNew) The category of the platform (e.g., `host`, `device`, `database`).
* `type` - (Required, ForceNew) The type within the category (e.g., `linux`, `windows`).
* `charset` - (Optional) The character set (e.g., `utf-8`, `gbk`). Default: `utf-8`.
* `comment` - (Optional) A description for the platform.
* `protocols` - (Optional) JSON string of supported protocols.
* `automation` - (Optional) JSON string of automation configuration.

## Attribute Reference

* `id` - The ID of the platform (numeric).
