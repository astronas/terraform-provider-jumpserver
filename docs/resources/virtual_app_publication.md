---
page_title: "jumpserver_virtual_app_publication Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a virtual app publication in JumpServer.
---

# `jumpserver_virtual_app_publication` Resource

The `jumpserver_virtual_app_publication` resource allows you to publish virtual applications to specific provider hosts.

## Example Usage

```hcl
resource "jumpserver_virtual_app_publication" "vscode_pub" {
  app           = "550e8400-e29b-41d4-a716-446655440000"
  provider_host = "660e8400-e29b-41d4-a716-446655440001"
  status        = "pending"
  comment       = "Publish VS Code to provider host"
}
```

## Argument Reference

* `app` - (Required, ForceNew) The UUID of the virtual app.
* `provider_host` - (Required, ForceNew) The UUID of the app-provider host.
* `status` - (Optional) The publication status. Default: `pending`.
* `comment` - (Optional) A description for the publication.

## Attribute Reference

* `id` - The UUID of the virtual app publication.
