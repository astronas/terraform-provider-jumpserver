<p align="center">
  <img src="https://download.jumpserver.org/images/jumpserver-logo.svg" alt="JumpServer Logo" width="400">
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
      version = "~> 1.5.0"
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
| `jumpserver_user_group_relation` | Manage user-to-group relationships | [user_group_relation.md](docs/resources/user_group_relation.md) |

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
| `jumpserver_directory` | Manage directory assets | [directory.md](docs/resources/directory.md) |
| `jumpserver_gpt` | Manage GPT assets | [gpt.md](docs/resources/gpt.md) |
| `jumpserver_asset` | Legacy assets (hostname/IP) | [asset.md](docs/resources/asset.md) |
| `jumpserver_node` | Manage asset tree nodes | [node.md](docs/resources/node.md) |
| `jumpserver_zone` | Manage network zones | [zone.md](docs/resources/zone.md) |
| `jumpserver_label` | Manage asset labels | [label.md](docs/resources/label.md) |
| `jumpserver_labeled_resource` | Attach labels to resources | [labeled_resource.md](docs/resources/labeled_resource.md) |
| `jumpserver_platform` | Manage platforms | [platform.md](docs/resources/platform.md) |
| `jumpserver_protocol_setting` | Manage protocol settings | [protocol_setting.md](docs/resources/protocol_setting.md) |
| `jumpserver_favorite_asset` | Mark an asset as favorite | [favorite_asset.md](docs/resources/favorite_asset.md) |

### Accounts

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_account` | Manage asset accounts | [account.md](docs/resources/account.md) |
| `jumpserver_account_template` | Manage account templates | [account_template.md](docs/resources/account_template.md) |
| `jumpserver_account_backup_plan` | Manage account backup plans | [account_backup_plan.md](docs/resources/account_backup_plan.md) |
| `jumpserver_virtual_account` | Manage virtual accounts | [virtual_account.md](docs/resources/virtual_account.md) |
| `jumpserver_account_risk` | Manage account risk entries | [account_risk.md](docs/resources/account_risk.md) |
| `jumpserver_gathered_account` | Manage gathered (discovered) accounts | [gathered_account.md](docs/resources/gathered_account.md) |
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
| `jumpserver_asset_perm_user_relation` | Link user to asset permission | [asset_perm_user_relation.md](docs/resources/asset_perm_user_relation.md) |
| `jumpserver_asset_perm_user_group_relation` | Link user group to asset permission | [asset_perm_user_group_relation.md](docs/resources/asset_perm_user_group_relation.md) |
| `jumpserver_asset_perm_asset_relation` | Link asset to asset permission | [asset_perm_asset_relation.md](docs/resources/asset_perm_asset_relation.md) |
| `jumpserver_asset_perm_node_relation` | Link node to asset permission | [asset_perm_node_relation.md](docs/resources/asset_perm_node_relation.md) |

### RBAC

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_org` | Manage organizations | [org.md](docs/resources/org.md) |
| `jumpserver_org_role` | Manage organization roles | [org_role.md](docs/resources/org_role.md) |
| `jumpserver_system_role` | Manage system roles | [system_role.md](docs/resources/system_role.md) |
| `jumpserver_rbac_role` | Manage RBAC roles | [rbac_role.md](docs/resources/rbac_role.md) |
| `jumpserver_role_binding` | Bind roles to users | [role_binding.md](docs/resources/role_binding.md) |
| `jumpserver_org_role_binding` | Bind org roles to users | [org_role_binding.md](docs/resources/org_role_binding.md) |
| `jumpserver_system_role_binding` | Bind system roles to users | [system_role_binding.md](docs/resources/system_role_binding.md) |

### Operations

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_job` | Manage automation jobs | [job.md](docs/resources/job.md) |
| `jumpserver_playbook` | Manage playbooks | [playbook.md](docs/resources/playbook.md) |
| `jumpserver_adhoc` | Manage ad-hoc commands | [adhoc.md](docs/resources/adhoc.md) |
| `jumpserver_ops_variable` | Manage ops variables | [ops_variable.md](docs/resources/ops_variable.md) |
| `jumpserver_ticket_flow` | Manage ticket approval flows | [ticket_flow.md](docs/resources/ticket_flow.md) |
| `jumpserver_ticket` | Manage tickets | [ticket.md](docs/resources/ticket.md) |
| `jumpserver_apply_asset_ticket` | Manage apply-asset tickets | [apply_asset_ticket.md](docs/resources/apply_asset_ticket.md) |
| `jumpserver_apply_command_ticket` | Manage apply-command tickets | [apply_command_ticket.md](docs/resources/apply_command_ticket.md) |
| `jumpserver_apply_login_ticket` | Manage apply-login tickets | [apply_login_ticket.md](docs/resources/apply_login_ticket.md) |
| `jumpserver_apply_login_asset_ticket` | Manage apply-login-asset tickets | [apply_login_asset_ticket.md](docs/resources/apply_login_asset_ticket.md) |
| `jumpserver_ticket_comment` | Manage ticket comments | [ticket_comment.md](docs/resources/ticket_comment.md) |

### Infrastructure

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_endpoint` | Manage terminal endpoints | [endpoint.md](docs/resources/endpoint.md) |
| `jumpserver_endpoint_rule` | Manage endpoint rules | [endpoint_rule.md](docs/resources/endpoint_rule.md) |
| `jumpserver_command_storage` | Manage command storage backends | [command_storage.md](docs/resources/command_storage.md) |
| `jumpserver_replay_storage` | Manage replay storage backends | [replay_storage.md](docs/resources/replay_storage.md) |
| `jumpserver_integration_application` | Manage integration applications | [integration_application.md](docs/resources/integration_application.md) |
| `jumpserver_applet` | Manage applets | [applet.md](docs/resources/applet.md) |
| `jumpserver_applet_host` | Manage applet hosts | [applet_host.md](docs/resources/applet_host.md) |
| `jumpserver_applet_host_deployment` | Manage applet host deployments | [applet_host_deployment.md](docs/resources/applet_host_deployment.md) |
| `jumpserver_applet_publication` | Manage applet publications | [applet_publication.md](docs/resources/applet_publication.md) |
| `jumpserver_virtual_app` | Manage virtual apps | [virtual_app.md](docs/resources/virtual_app.md) |
| `jumpserver_virtual_app_publication` | Manage virtual app publications | [virtual_app_publication.md](docs/resources/virtual_app_publication.md) |
| `jumpserver_app_provider` | Manage application providers | [app_provider.md](docs/resources/app_provider.md) |

### Terminal & Sessions

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_terminal` | Manage terminals | [terminal.md](docs/resources/terminal.md) |
| `jumpserver_terminal_registration` | Register new terminals | [terminal_registration.md](docs/resources/terminal_registration.md) |
| `jumpserver_session` | Manage sessions | [session.md](docs/resources/session.md) |
| `jumpserver_session_sharing` | Manage session sharing | [session_sharing.md](docs/resources/session_sharing.md) |
| `jumpserver_session_join_record` | Manage session join records | [session_join_record.md](docs/resources/session_join_record.md) |
| `jumpserver_terminal_command` | Manage terminal command records | [terminal_command.md](docs/resources/terminal_command.md) |

### Authentication

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_access_key` | Manage API access keys | [access_key.md](docs/resources/access_key.md) |
| `jumpserver_ssh_key` | Manage SSH keys | [ssh_key.md](docs/resources/ssh_key.md) |
| `jumpserver_passkey` | Manage passkeys (WebAuthn) | [passkey.md](docs/resources/passkey.md) |
| `jumpserver_connection_token` | Manage connection tokens | [connection_token.md](docs/resources/connection_token.md) |
| `jumpserver_super_connection_token` | Manage super connection tokens | [super_connection_token.md](docs/resources/super_connection_token.md) |
| `jumpserver_temp_token` | Manage temporary tokens | [temp_token.md](docs/resources/temp_token.md) |

### Settings

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_setting` | Manage basic settings (singleton) | [setting.md](docs/resources/setting.md) |
| `jumpserver_leak_password` | Manage leaked password entries | [leak_password.md](docs/resources/leak_password.md) |
| `jumpserver_chatai_prompt` | Manage ChatAI prompts | [chatai_prompt.md](docs/resources/chatai_prompt.md) |

## Data Sources

| Data Source | Description | Docs |
|---|---|---|
| `jumpserver_user` | Look up a user by username | [user.md](docs/data-sources/user.md) |
| `jumpserver_user_group` | Look up a user group by name | [user_group.md](docs/data-sources/user_group.md) |
| `jumpserver_node` | Look up a node by value | [node.md](docs/data-sources/node.md) |
| `jumpserver_zone` | Look up a zone by name | [zone.md](docs/data-sources/zone.md) |
| `jumpserver_platform` | Look up a platform by name | [platform.md](docs/data-sources/platform.md) |
| `jumpserver_asset` | Look up an asset by name | [asset.md](docs/data-sources/asset.md) |
| `jumpserver_org` | Look up an organization by name | [org.md](docs/data-sources/org.md) |
| `jumpserver_org_role` | Look up an org role by name | [org_role.md](docs/data-sources/org_role.md) |
| `jumpserver_system_role` | Look up a system role by name | [system_role.md](docs/data-sources/system_role.md) |
| `jumpserver_host` | Look up a host by name | [host.md](docs/data-sources/host.md) |
| `jumpserver_label` | Look up a label by name | [label.md](docs/data-sources/label.md) |
| `jumpserver_account` | Look up an account by name | [account.md](docs/data-sources/account.md) |
| `jumpserver_account_template` | Look up an account template by name | [account_template.md](docs/data-sources/account_template.md) |
| `jumpserver_endpoint` | Look up an endpoint by name | [endpoint.md](docs/data-sources/endpoint.md) |
| `jumpserver_database` | Look up a database asset by name | [database.md](docs/data-sources/database.md) |
| `jumpserver_gateway` | Look up a gateway by name | [gateway.md](docs/data-sources/gateway.md) |
| `jumpserver_applet` | Look up an applet by name | [applet.md](docs/data-sources/applet.md) |
| `jumpserver_device` | Look up a device by name | [device.md](docs/data-sources/device.md) |
| `jumpserver_web` | Look up a web asset by name | [web.md](docs/data-sources/web.md) |
| `jumpserver_cloud` | Look up a cloud asset by name | [cloud.md](docs/data-sources/cloud.md) |
| `jumpserver_custom` | Look up a custom asset by name | [custom.md](docs/data-sources/custom.md) |
| `jumpserver_asset_permission` | Look up an asset permission by name | [asset_permission.md](docs/data-sources/asset_permission.md) |
| `jumpserver_asset_category` | Look up an asset category by name | [asset_category.md](docs/data-sources/asset_category.md) |
| `jumpserver_protocol` | Look up a protocol by name | [protocol.md](docs/data-sources/protocol.md) |
| `jumpserver_content_type` | Look up a content type by name | [content_type.md](docs/data-sources/content_type.md) |
| `jumpserver_terminal` | Look up a terminal by name | [terminal.md](docs/data-sources/terminal.md) |
| `jumpserver_ticket_flow` | Look up a ticket flow by name | [ticket_flow.md](docs/data-sources/ticket_flow.md) |
| `jumpserver_login_log` | Look up login logs by username | [login_log.md](docs/data-sources/login_log.md) |
| `jumpserver_operate_log` | Look up operate logs by user | [operate_log.md](docs/data-sources/operate_log.md) |
| `jumpserver_activity_log` | Look up activity logs by resource ID | [activity_log.md](docs/data-sources/activity_log.md) |
| `jumpserver_rbac_permission` | Look up an RBAC permission by name | [rbac_permission.md](docs/data-sources/rbac_permission.md) |
| `jumpserver_server_info` | Retrieve server information | [server_info.md](docs/data-sources/server_info.md) |
| `jumpserver_site_message` | Look up a site message by subject | [site_message.md](docs/data-sources/site_message.md) |
| `jumpserver_sms_backend` | Look up an SMS backend by name | [sms_backend.md](docs/data-sources/sms_backend.md) |
| `jumpserver_ops_task` | Look up an ops task by name | [ops_task.md](docs/data-sources/ops_task.md) |
| `jumpserver_user_msg_subscription` | Look up user message subscriptions | [user_msg_subscription.md](docs/data-sources/user_msg_subscription.md) |

### Notifications

| Resource | Description | Docs |
|---|---|---|
| `jumpserver_system_msg_subscription` | Manage system message subscriptions | [system_msg_subscription.md](docs/resources/system_msg_subscription.md) |

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
