---
page_title: "jumpserver_server_info Data Source - terraform-provider-jumpserver"
subcategory: "Settings"
description: |-
  Retrieves JumpServer server information.
---

# `jumpserver_server_info` Data Source

Use this data source to retrieve the current server information from JumpServer.

## Example Usage

```hcl
data "jumpserver_server_info" "current" {}
```

## Argument Reference

No arguments are required.

## Attribute Reference

* `id` - Always `server-info`.
* `current_time` - The server's current time.
