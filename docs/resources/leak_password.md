---
page_title: "jumpserver_leak_password Resource - terraform-provider-jumpserver"
subcategory: "Settings"
description: |-
  Manages a leaked password entry in JumpServer.
---

# `jumpserver_leak_password` Resource

Manages leaked password entries in JumpServer for security monitoring.

## Example Usage

```hcl
resource "jumpserver_leak_password" "example" {
  password = "compromised_password_hash"
}
```

## Argument Reference

- **`password`** - (Required, Sensitive) The leaked password value to register.

## Attribute Reference

- **`id`** - The ID of the leak password entry.
