---
page_title: "jumpserver_endpoint Resource - terraform-provider-jumpserver"
subcategory: "Infrastructure"
description: |-
  Manages a JumpServer Endpoint.
---

# jumpserver_endpoint

Manages a JumpServer Endpoint. Endpoints define network access points with protocol-specific ports for SSH, RDP, database, and other connections.

## Example Usage

```hcl
resource "jumpserver_endpoint" "main" {
  name            = "main-endpoint"
  host            = "jumpserver.example.com"
  https_port      = 443
  ssh_port        = 2222
  rdp_port        = 3389
  mysql_port      = 33061
  postgresql_port = 54320
  is_active       = true
  comment         = "Primary access endpoint"
}

resource "jumpserver_endpoint" "internal" {
  name     = "internal-endpoint"
  host     = "10.0.1.100"
  ssh_port = 22
  rdp_port = 3389
  comment  = "Internal network endpoint"
}
```

## Argument Reference

- `name` - (Required) The name of the endpoint.
- `host` - (Optional) Access address. Empty uses the browser address.
- `https_port` - (Optional) HTTPS port. Defaults to `443`.
- `http_port` - (Optional) HTTP port. Defaults to `80`.
- `ssh_port` - (Optional) SSH port. Defaults to `2222`.
- `rdp_port` - (Optional) RDP port. Defaults to `3389`.
- `mysql_port` - (Optional) MySQL port. Defaults to `33061`.
- `mariadb_port` - (Optional) MariaDB port. Defaults to `33062`.
- `postgresql_port` - (Optional) PostgreSQL port. Defaults to `54320`.
- `redis_port` - (Optional) Redis port. Defaults to `63790`.
- `vnc_port` - (Optional) VNC port. Defaults to `15900`.
- `oracle_port` - (Optional) Oracle port. Defaults to `15210`.
- `sqlserver_port` - (Optional) SQL Server port. Defaults to `14330`.
- `mongodb_port` - (Optional) MongoDB port. Defaults to `27018`.
- `is_active` - (Optional) Whether the endpoint is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the endpoint.

## Import

```shell
terraform import jumpserver_endpoint.example <uuid>
```
