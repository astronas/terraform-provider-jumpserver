---
page_title: "JumpServer Provider"
description: |-
  The JumpServer provider allows you to manage JumpServer resources using Terraform.
---

# JumpServer Provider

The JumpServer provider allows you to manage [JumpServer](https://www.jumpserver.org/) resources using Terraform — including users, assets, accounts, permissions, RBAC, operations, and infrastructure.

## Example Usage

```hcl
provider "jumpserver" {
  base_url        = "https://jumpserver.example.com"
  username        = "admin"
  password        = "adminpass"
  skip_tls_verify = false   # optional, default: false
  api_version     = "v1"    # optional, "v1" or "v2", default: "v1"
}
```

## Authentication

The provider authenticates using username/password to obtain a Bearer token from the JumpServer API.

## Argument Reference

* `base_url` (Required) — The base URL of your JumpServer instance.
* `username` (Required) — The username used to authenticate with JumpServer.
* `password` (Required, Sensitive) — The password used to authenticate with JumpServer.
* `skip_tls_verify` (Optional) — Skip TLS certificate validation. Default: `false`.
* `api_version` (Optional) — API version to use (`v1` or `v2`). Default: `v1`.

## Resources

### Users & Groups
- `jumpserver_user` — Manage users
- `jumpserver_user_group` — Manage user groups
- `jumpserver_user_group_relation` — Manage user-to-group relationships

### Assets
- `jumpserver_host` — Manage host assets
- `jumpserver_device` — Manage device assets
- `jumpserver_database` — Manage database assets
- `jumpserver_web` — Manage web assets
- `jumpserver_cloud` — Manage cloud assets
- `jumpserver_gateway` — Manage gateways
- `jumpserver_custom` — Manage custom assets
- `jumpserver_directory` — Manage directory assets
- `jumpserver_gpt` — Manage GPT assets
- `jumpserver_asset` — Manage generic assets
- `jumpserver_node` — Manage nodes
- `jumpserver_zone` — Manage zones
- `jumpserver_label` — Manage labels
- `jumpserver_labeled_resource` — Attach labels to resources
- `jumpserver_platform` — Manage platforms

### Assets (continued)
- `jumpserver_protocol_setting` — Manage protocol settings

### Accounts
- `jumpserver_account` — Manage accounts
- `jumpserver_account_template` — Manage account templates
- `jumpserver_account_backup_plan` — Manage account backup plans
- `jumpserver_virtual_account` — Manage virtual accounts
- `jumpserver_account_risk` — Manage account risk entries
- `jumpserver_system_user` — Manage system users (legacy)

### Permissions & ACLs
- `jumpserver_asset_permission` — Manage asset permissions
- `jumpserver_command_group` — Manage command groups
- `jumpserver_command_filter_acl` — Manage command filter ACLs
- `jumpserver_login_acl` — Manage login ACLs
- `jumpserver_login_asset_acl` — Manage login asset ACLs
- `jumpserver_connect_method_acl` — Manage connect method ACLs
- `jumpserver_data_masking_rule` — Manage data masking rules

### RBAC
- `jumpserver_org` — Manage organizations
- `jumpserver_org_role` — Manage organization roles
- `jumpserver_system_role` — Manage system roles
- `jumpserver_role_binding` — Manage role bindings
- `jumpserver_org_role_binding` — Manage organization role bindings
- `jumpserver_system_role_binding` — Manage system role bindings

### Ops
- `jumpserver_job` — Manage jobs
- `jumpserver_playbook` — Manage playbooks
- `jumpserver_adhoc` — Manage ad-hoc commands
- `jumpserver_ops_variable` — Manage ops variables
- `jumpserver_ticket_flow` — Manage ticket approval flows

### Terminal
- `jumpserver_endpoint` — Manage endpoints
- `jumpserver_endpoint_rule` — Manage endpoint rules
- `jumpserver_command_storage` — Manage command storages
- `jumpserver_replay_storage` — Manage replay storages
- `jumpserver_applet` — Manage applets
- `jumpserver_applet_host` — Manage applet hosts
- `jumpserver_applet_host_deployment` — Manage applet host deployments
- `jumpserver_applet_publication` — Manage applet publications
- `jumpserver_virtual_app` — Manage virtual apps
- `jumpserver_virtual_app_publication` — Manage virtual app publications
- `jumpserver_app_provider` — Manage application providers

### Authentication
- `jumpserver_access_key` — Manage API access keys
- `jumpserver_integration_application` — Manage integration applications
- `jumpserver_ssh_key` — Manage SSH keys
- `jumpserver_passkey` — Manage passkeys (WebAuthn)

### Settings
- `jumpserver_leak_password` — Manage leaked password entries

## Data Sources

- `jumpserver_user` — Look up a user by username
- `jumpserver_user_group` — Look up a user group by name
- `jumpserver_node` — Look up a node by value
- `jumpserver_zone` — Look up a zone by name
- `jumpserver_platform` — Look up a platform by name
- `jumpserver_asset` — Look up an asset by name
- `jumpserver_org` — Look up an organization by name
- `jumpserver_org_role` — Look up an organization role by name
- `jumpserver_system_role` — Look up a system role by name
- `jumpserver_host` — Look up a host by name
- `jumpserver_label` — Look up a label by name
- `jumpserver_account` — Look up an account by name
- `jumpserver_account_template` — Look up an account template by name
- `jumpserver_endpoint` — Look up an endpoint by name
- `jumpserver_database` — Look up a database asset by name
- `jumpserver_gateway` — Look up a gateway by name
- `jumpserver_applet` — Look up an applet by name