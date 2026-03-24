---
page_title: "jumpserver_database Resource"
subcategory: "Assets"
description: |-
  Manages a database asset in JumpServer.
---

# jumpserver_database

Manages a database asset in JumpServer. Supports PostgreSQL, MySQL, Oracle, MongoDB and other database types with SSL/TLS configuration.

## Example Usage

```hcl
resource "jumpserver_database" "postgres_main" {
  name      = "postgres-main"
  address   = "db.example.com"
  platform  = 5
  zone_name = "Production"
  node_name = "Databases"
  db_name   = "myapp"
  use_ssl   = true
  pg_ssl_mode = "require"
  comment   = "Main PostgreSQL database"

  protocols {
    name = "postgresql"
    port = 5432
  }

  accounts {
    name        = "dbadmin"
    username    = "dbadmin"
    secret_type = "password"
    secret      = "secure-password"
  }
}
```

## Argument Reference

### Base fields (shared with all asset types)
- `name` - (Required) The name of the asset.
- `address` - (Required) The address (IP/hostname) of the database server.
- `platform` - (Required) The platform ID (integer).
- `zone_name` - (Required) The name of the zone.
- `node_name` - (Required) The name of the node.
- `comment` - (Optional) A description.
- `is_active` - (Optional) Whether the asset is active. Default: `true`.
- `accounts` - (Optional) List of account blocks. See Host resource for schema.
- `protocols` - (Optional) List of protocol blocks (`name`, `port`).

### Database-specific fields
- `db_name` - (Required) The database name.
- `use_ssl` - (Optional) Use SSL connection. Default: `false`.
- `allow_invalid_cert` - (Optional) Allow invalid SSL certificates. Default: `false`.
- `ca_cert` - (Optional, Sensitive) CA certificate for SSL.
- `client_cert` - (Optional, Sensitive) Client certificate for SSL.
- `client_key` - (Optional, Sensitive) Client key for SSL.
- `pg_ssl_mode` - (Optional) PostgreSQL SSL mode: `prefer`, `require`, `disable`, `allow`. Default: `prefer`.

## Attribute Reference

- `id` - The UUID of the database asset.
- `zone_id` - The resolved zone UUID.
- `node_ids` - The resolved node UUIDs.
