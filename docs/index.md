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