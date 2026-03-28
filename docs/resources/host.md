---
page_title: "jumpserver_host Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a host asset in JumpServer.
---

# `jumpserver_host` Resource

The `jumpserver_host` resource allows you to create and manage *hosts* in Jumpserver. A host represents a specific server/endpoint that Jumpserver will manage. This resource also supports creating SSH accounts (or other protocols) on the host.

## Example Usage

```hcl
resource "jumpserver_host" "example_host" {
  # Basic host info
  name     = "server-lxc1"
  address  = "10.10.10.50"
  platform = 32
  comment  = "Production Linux server"

  # Zone and Node by name
  zone_name = "Production"
  node_name   = "Linux Servers"

  # Define SSH accounts
  accounts {
    on_invalid  = "error"
    is_active   = true
    name        = "root"
    username    = "root"
    secret_type = "ssh_key"
    secret      = file("${path.module}/ssh_key/id_ed25519")
  }

  # Define protocols
  protocols {
    name = "ssh"
    port = 22
  }
  protocols {
    name = "sftp"
    port = 22
  }
}
```

## Argument Reference

- **`name`** - (Required) The name of the host in Jumpserver.
- **`address`** - (Required) The IP address (or hostname) of the host.
- **`platform`** - (Required) The platform code for this host (e.g., `32` for Linux).
- **`comment`** - (Optional) A comment or description for the host, you can search host by comment in jumpserver.

- **`zone_name`** - (Required) The **name** of the Zone in Jumpserver that this host should belong to. The provider will look up the Zone by its `name` and retrieve its ID to associate the host.
- **`node_name`** - (Required) The **name** of the Node in Jumpserver that this host should belong to. The provider will look up the Node by its `name` and retrieve its ID to associate the host.

- **`accounts`** - (Optional) A list of account definitions for this host.
    - **`on_invalid`** - (Optional) Action if the credential becomes invalid. Defaults to `"error"`.
    - **`is_active`** - (Optional) Whether the account is active. Defaults to `true`.
    - **`name`** - (Required) An identifier for the account (e.g., `"root"`).
    - **`username`** - (Required) The actual username on the host.
    - **`secret_type`** - (Required) The type of secret (e.g., `"ssh_key"` or `"password"`).
    - **`secret`** - (Required, Sensitive) The key or password used for authentication.

- **`protocols`** - (Optional) A list of protocols the host can be accessed by.
    - **`name`** - (Required) The name of the protocol (e.g., `"ssh"`, `"sftp"`).
    - **`port`** - (Required) The port number for that protocol.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- **`id`** - The ID of the host in Jumpserver.
- **`zone_id`** - The actual zone ID used in Jumpserver. (Computed at create-time if you supply `zone_name`.)
- **`node_ids`** - A list of node IDs (in Jumpserver) this host is attached to. (Computed at create-time if you supply `node_name`.)

## Notes

- If the specified `zone_name` or `node_name` do not exist in Jumpserver, creation of the host fails. This resource does **not** create or delete zones/nodes.

## Import

Hosts can be imported using their JumpServer UUID:

```shell
terraform import jumpserver_host.example_host <host-uuid>
```

On import, `zone_name` and `node_name` are automatically resolved from the API.
- During updates:
    - If you change `zone_name` or `node_name`, the provider will look up new IDs and update the host accordingly.
- During `destroy`, only the host is deleted. Zones and nodes remain intact.
