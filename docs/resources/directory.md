---
page_title: "jumpserver_directory Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a directory asset in JumpServer.
---

# `jumpserver_directory` Resource

The `jumpserver_directory` resource allows you to create and manage directory-type assets in JumpServer such as NFS, CIFS, or FTP shares.

## Example Usage

```hcl
resource "jumpserver_directory" "shared" {
  name      = "shared-storage"
  address   = "10.10.10.100"
  platform  = 40
  comment   = "Shared NFS storage"
  zone_name = "Production"
  node_name = "Storage"

  protocols {
    name = "ftp"
    port = 21
  }
}
```

## Argument Reference

* `name` - (Required) The name of the directory asset.
* `address` - (Required) The IP address or hostname of the directory.
* `platform` - (Required) The platform code for this directory.
* `comment` - (Optional) A description for the directory.
* `zone_name` - (Required) The name of the Zone this directory belongs to.
* `node_name` - (Required) The name of the Node this directory belongs to.
* `accounts` - (Optional) A list of account definitions. See the `jumpserver_host` resource for sub-attributes.
* `protocols` - (Optional) A list of protocols. See the `jumpserver_host` resource for sub-attributes.

## Attribute Reference

* `id` - The UUID of the directory asset.
