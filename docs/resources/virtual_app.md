---
page_title: "jumpserver_virtual_app Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a virtual app in JumpServer.
---

# `jumpserver_virtual_app` Resource

The `jumpserver_virtual_app` resource allows you to create and manage virtual applications (Docker-based apps) in JumpServer.

## Example Usage

```hcl
resource "jumpserver_virtual_app" "vscode" {
  name           = "vscode"
  display_name   = "VS Code"
  image_name     = "jumpserver/vscode:latest"
  version        = "1.0.0"
  author         = "admin"
  is_active      = true
  image_protocol = "vnc"
  image_port     = 5900
  comment        = "VS Code virtual app"
}
```

## Argument Reference

* `name` - (Required) The name of the virtual app.
* `display_name` - (Required) The display name.
* `image_name` - (Required) The Docker image name.
* `version` - (Required) The version.
* `author` - (Optional) The author.
* `is_active` - (Optional) Whether the virtual app is active. Default: `true`.
* `image_protocol` - (Optional) The protocol used by the image (`vnc`, `rdp`). Default: `vnc`.
* `image_port` - (Optional) The port used by the image. Default: `5900`.
* `protocols` - (Optional) JSON string of supported protocols.
* `tags` - (Optional) JSON string of tags.
* `comment` - (Optional) A description for the virtual app.

## Attribute Reference

* `id` - The UUID of the virtual app.
