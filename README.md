<p align="center">
  <img src="https://download.jumpserver.org/images/jumpserver-logo.svg" alt="JumpServer Logo" width="200">
</p>

<h1 align="center">Terraform Provider for JumpServer</h1>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
</p>

Manage [JumpServer](https://www.jumpserver.org/) resources with Terraform. JumpServer is an open-source bastion host for managing assets, users, permissions, and operations.

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "astronas/jumpserver"
      version = "~> 1.0.0"
    }
  }
}

provider "jumpserver" {
  base_url        = "https://jumpserver.example.com"
  username        = "admin"
  password        = "adminpass"
  skip_tls_verify = false   # optional, default: false
  api_version     = "v1"    # optional, "v1" or "v2", default: "v1"
}
```

## Resources

### Users & Groups

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_user` | Manage users | [user.md](docs/resources/user.md) |
| `jumpserver_user_group` | Manage user groups | [user_group.md](docs/resources/user_group.md) |

### Assets

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_host` | Manage host assets | [host.md](docs/resources/host.md) |
| `jumpserver_database` | Manage database assets | [database.md](docs/resources/database.md) |
| `jumpserver_device` | Manage network device assets | [device.md](docs/resources/device.md) |
| `jumpserver_web` | Manage web assets | [web.md](docs/resources/web.md) |
| `jumpserver_gateway` | Manage gateway assets | [gateway.md](docs/resources/gateway.md) |
| `jumpserver_cloud` | Manage cloud assets | [cloud.md](docs/resources/cloud.md) |
| `jumpserver_custom` | Manage custom assets | [custom.md](docs/resources/custom.md) |
| `jumpserver_asset` | Legacy assets (hostname/IP) | [asset.md](docs/resources/asset.md) |
| `jumpserver_node` | Manage asset tree nodes | [node.md](docs/resources/node.md) |
| `jumpserver_zone` | Manage network zones | [zone.md](docs/resources/zone.md) |
| `jumpserver_label` | Manage asset labels | [label.md](docs/resources/label.md) |

### Accounts

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_account` | Manage asset accounts | [account.md](docs/resources/account.md) |
| `jumpserver_account_template` | Manage account templates | [account_template.md](docs/resources/account_template.md) |
| `jumpserver_account_backup_plan` | Manage account backup plans | [account_backup_plan.md](docs/resources/account_backup_plan.md) |
| `jumpserver_system_user` | Legacy system users | [system_user.md](docs/resources/system_user.md) |

### Permissions & ACLs

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_asset_permission` | Manage asset access permissions | [asset_permission.md](docs/resources/asset_permission.md) |
| `jumpserver_command_group` | Manage command groups | [command_group.md](docs/resources/command_group.md) |
| `jumpserver_command_filter_acl` | Manage command filter ACLs | [command_filter_acl.md](docs/resources/command_filter_acl.md) |
| `jumpserver_login_acl` | Manage login ACLs | [login_acl.md](docs/resources/login_acl.md) |
| `jumpserver_login_asset_acl` | Manage login asset ACLs | [login_asset_acl.md](docs/resources/login_asset_acl.md) |
| `jumpserver_connect_method_acl` | Manage connect method ACLs | [connect_method_acl.md](docs/resources/connect_method_acl.md) |
| `jumpserver_data_masking_rule` | Manage data masking rules | [data_masking_rule.md](docs/resources/data_masking_rule.md) |

### RBAC

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_org_role` | Manage organization roles | [org_role.md](docs/resources/org_role.md) |
| `jumpserver_system_role` | Manage system roles | [system_role.md](docs/resources/system_role.md) |
| `jumpserver_role_binding` | Bind roles to users | [role_binding.md](docs/resources/role_binding.md) |

### Operations

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_job` | Manage automation jobs | [job.md](docs/resources/job.md) |
| `jumpserver_playbook` | Manage playbooks | [playbook.md](docs/resources/playbook.md) |

### Infrastructure

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_endpoint` | Manage terminal endpoints | [endpoint.md](docs/resources/endpoint.md) |
| `jumpserver_endpoint_rule` | Manage endpoint rules | [endpoint_rule.md](docs/resources/endpoint_rule.md) |
| `jumpserver_command_storage` | Manage command storage backends | [command_storage.md](docs/resources/command_storage.md) |
| `jumpserver_replay_storage` | Manage replay storage backends | [replay_storage.md](docs/resources/replay_storage.md) |
| `jumpserver_integration_application` | Manage integration applications | [integration_application.md](docs/resources/integration_application.md) |

## Example Usage

```hcl
resource "jumpserver_user_group" "ops" {
  name    = "ops-team"
  comment = "Operations team"
}

resource "jumpserver_user" "alice" {
  name         = "Alice"
  username     = "alice"
  email        = "alice@example.com"
  is_active    = true
  system_roles = ["User"]
}

resource "jumpserver_node" "prod" {
  value = "Production"
}

resource "jumpserver_host" "web01" {
  name      = "web01"
  address   = "192.168.1.10"
  platform  = 1
  zone_name = "Default"
  node_name = "Production"

  protocols {
    name = "ssh"
    port = 22
  }

  accounts {
    name        = "root"
    username    = "root"
    secret_type = "password"
    secret      = "s3cr3t"
  }
}

resource "jumpserver_asset_permission" "allow_ops" {
  name         = "allow-ops-to-prod"
  is_active    = true
  user_groups  = [jumpserver_user_group.ops.id]
  assets       = [jumpserver_host.web01.id]
  system_users = []
  actions      = ["connect"]
}
```

## License

MIT License — see [LICENSE](LICENSE) for details.

## Author

**Thibaut Gianola** ([@astronas](https://github.com/astronas))
