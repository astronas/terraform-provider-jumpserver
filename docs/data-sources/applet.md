---
page_title: "jumpserver_applet Data Source - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Look up an applet by name in JumpServer.
---

# `jumpserver_applet` Data Source

Use this data source to look up an existing applet by name.

## Example Usage

```hcl
data "jumpserver_applet" "chrome" {
  name = "chrome"
}
```

## Argument Reference

- **`name`** - (Required) The name of the applet to look up.

## Attribute Reference

- **`id`** - The UUID of the applet.
- **`display_name`** - The display name of the applet.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the applet is active.
