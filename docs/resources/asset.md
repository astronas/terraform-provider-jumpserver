---
page_title: "jumpserver_asset Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a legacy asset in JumpServer.
---

# `jumpserver_asset` Resource

The jumpserver_asset resource allows you to create and manage assets in Jumpserver. An asset represents a device or server that you want to manage using Jumpserver.


## Example Usage

```hcl
resource "jumpserver_asset" "example_asset" {
  hostname      = "server-vm1"
  ip            = "X.X.X.X"
  platform      = "Linux"
  protocols     = ["ssh/22"]
  nodes_display = ["/Default/NODE1"]
}
```

## Argument Reference

* `hostname` - (Required) The hostname of the asset.
* `ip` - (Required) The IP address of the asset.
* `platform` - (Required) The platform of the asset (e.g., Linux).
* `protocols` - (Optional) List of protocols the asset supports.
* `nodes_display` - (Optional) List of nodes the asset is associated with.

## Attribute Reference

* `id` - The ID of the asset.
* `hostname` - The hostname of the asset.
* `ip` - The IP address of the asset.
* `platform` - The platform of the asset.
* `protocols` - List of protocols the asset supports.
* `nodes_display` - List of nodes the asset is associated with.