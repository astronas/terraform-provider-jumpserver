---
page_title: "jumpserver_applet_host Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages an applet host in JumpServer.
---

# `jumpserver_applet_host` Resource

The `jumpserver_applet_host` resource allows you to create and manage applet hosts (Windows RDP hosts that run applets) in JumpServer.

## Example Usage

```hcl
resource "jumpserver_applet_host" "rdp_host" {
  name     = "applet-host-01"
  address  = "10.10.10.50"
  platform = 5
  comment  = "Applet host for Chrome"

  zone_name = "Production"
  node_name = "Applet Hosts"

  accounts {
    name        = "admin"
    username    = "administrator"
    secret_type = "password"
    secret      = "SecurePassword123"
    is_active   = true
  }

  protocols {
    name = "rdp"
    port = 3389
  }
}
```

## Argument Reference

* `name` - (Required) The name of the applet host.
* `address` - (Required) The IP address or hostname.
* `platform` - (Required) The platform code.
* `comment` - (Optional) A description for the applet host.
* `zone_name` - (Optional) The name of the Zone.
* `node_name` - (Optional) The name of the Node.
* `accounts` - (Optional) A list of account definitions. See the `jumpserver_host` resource for sub-attributes.
* `protocols` - (Optional) A list of protocols. See the `jumpserver_host` resource for sub-attributes.
* `deploy_options` - (Optional) JSON string of deployment options.

## Attribute Reference

* `id` - The UUID of the applet host.
