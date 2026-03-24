---
page_title: "jumpserver_account_risk Resource - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Manages an account risk entry in JumpServer.
---

# `jumpserver_account_risk` Resource

Manages account risk entries in JumpServer for security monitoring and risk tracking.

## Example Usage

```hcl
resource "jumpserver_account_risk" "example" {
  username = "admin"
  asset    = jumpserver_host.example.id
  status   = "pending"
  details  = "Weak password detected"
}
```

## Argument Reference

- **`username`** - (Required) The username associated with the risk entry.
- **`asset`** - (Optional) The UUID of the asset associated with the risk.
- **`status`** - (Optional) The status of the risk (`pending`, `confirmed`, `ignored`). Defaults to `pending`.
- **`details`** - (Optional) Additional details about the account risk.

## Attribute Reference

- **`id`** - The UUID of the account risk entry.
- **`risk`** - The computed risk level (JSON).
- **`date_created`** - The creation date.
