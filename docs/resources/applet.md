---
page_title: "jumpserver_applet Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages an applet in JumpServer.
---

# `jumpserver_applet` Resource

The `jumpserver_applet` resource allows you to create and manage applets (remote applications) in JumpServer.

## Example Usage

```hcl
resource "jumpserver_applet" "chrome" {
  name         = "chrome"
  display_name = "Google Chrome"
  version      = "1.0.0"
  author       = "admin"
  type         = "general"
  is_active    = true
  comment      = "Chrome browser applet"
}
```

## Argument Reference

* `name` - (Required) The name of the applet.
* `display_name` - (Required) The display name of the applet.
* `version` - (Required) The version of the applet.
* `author` - (Optional) The author of the applet.
* `type` - (Optional) The type of the applet. Default: `general`.
* `is_active` - (Optional) Whether the applet is active. Default: `true`.
* `can_concurrent` - (Optional) Whether the applet supports concurrent sessions. Default: `false`.
* `protocols` - (Optional) JSON string of supported protocols.
* `tags` - (Optional) JSON string of tags.
* `comment` - (Optional) A description for the applet.
* `edition` - (Optional) The edition of the applet. Default: `community`.

## Attribute Reference

* `id` - The UUID of the applet.
