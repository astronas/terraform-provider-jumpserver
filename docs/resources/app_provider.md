---
page_title: "jumpserver_app_provider Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages an application provider in JumpServer.
---

# `jumpserver_app_provider` Resource

Manages application providers (terminal app providers) in JumpServer.

## Example Usage

```hcl
resource "jumpserver_app_provider" "example" {
  name     = "my-app-provider"
  hostname = "app-server-01.example.com"
}
```

## Argument Reference

- **`name`** - (Required) The name of the app provider.
- **`hostname`** - (Required) The hostname of the provider.
- **`terminal`** - (Optional) The UUID of the associated terminal.

## Attribute Reference

- **`id`** - The UUID of the app provider.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
