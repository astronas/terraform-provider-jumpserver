---
page_title: "jumpserver_passkey Resource - terraform-provider-jumpserver"
subcategory: "Authentication"
description: |-
  Manages a passkey (WebAuthn) in JumpServer.
---

# `jumpserver_passkey` Resource

Manages passkeys (WebAuthn/FIDO2 credentials) in JumpServer. Most fields are read-only as they are set during registration.

## Example Usage

```hcl
resource "jumpserver_passkey" "example" {
  is_active = true
}
```

## Argument Reference

- **`is_active`** - (Optional) Whether the passkey is active. Defaults to `true`.

## Attribute Reference

- **`id`** - The UUID of the passkey.
- **`name`** - The name of the passkey.
- **`platform`** - The platform that registered the passkey.
- **`created_by`** - Who created the passkey.
- **`date_last_used`** - The date the passkey was last used.
- **`date_created`** - The creation date.
