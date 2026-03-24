---
page_title: "jumpserver_applet_publication Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages an applet publication in JumpServer.
---

# `jumpserver_applet_publication` Resource

The `jumpserver_applet_publication` resource allows you to publish applets to specific applet hosts.

## Example Usage

```hcl
resource "jumpserver_applet_publication" "chrome_pub" {
  applet  = "550e8400-e29b-41d4-a716-446655440000"
  host    = "660e8400-e29b-41d4-a716-446655440001"
  status  = "pending"
  comment = "Publish Chrome to host-01"
}
```

## Argument Reference

* `applet` - (Required, ForceNew) The UUID of the applet to publish.
* `host` - (Required, ForceNew) The UUID of the applet host.
* `status` - (Optional) The publication status. Default: `pending`.
* `comment` - (Optional) A description for the publication.

## Attribute Reference

* `id` - The UUID of the applet publication.
