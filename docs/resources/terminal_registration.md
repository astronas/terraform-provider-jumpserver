---
page_title: "jumpserver_terminal_registration Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Registers a new terminal in JumpServer.
---

# `jumpserver_terminal_registration` Resource

The `jumpserver_terminal_registration` resource registers a new terminal component in JumpServer. This is a create-only resource — all fields force replacement on change.

## Example Usage

```hcl
resource "jumpserver_terminal_registration" "koko" {
  name    = "koko-01"
  type    = "koko"
  comment = "KoKo terminal for SSH"
}
```

## Argument Reference

- **`name`** - (Required, ForceNew) The name of the terminal.
- **`type`** - (Required, ForceNew) The type of the terminal.
- **`comment`** - (Optional, ForceNew) A comment or description.

## Attributes Reference

- **`id`** - The UUID of the registered terminal.
- **`service_account`** - The service account created for this terminal (JSON).
- **`remote_addr`** - The remote address.
