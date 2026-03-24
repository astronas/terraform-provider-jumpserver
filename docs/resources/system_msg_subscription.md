---
page_title: "jumpserver_system_msg_subscription Resource - terraform-provider-jumpserver"
subcategory: "Notifications"
description: |-
  Manages a system message subscription in JumpServer.
---

# `jumpserver_system_msg_subscription` Resource

The `jumpserver_system_msg_subscription` resource manages notification subscriptions for system messages in JumpServer. Subscriptions are identified by their message type.

## Example Usage

```hcl
resource "jumpserver_system_msg_subscription" "server_change" {
  message_type     = "server_account_change"
  users            = ["xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
  groups           = []
  receive_backends = ["site_msg", "email"]
}
```

## Argument Reference

- **`message_type`** - (Required, ForceNew) The message type identifier.
- **`users`** - (Optional) List of user UUIDs to subscribe.
- **`groups`** - (Optional) List of user group UUIDs to subscribe.
- **`receive_backends`** - (Optional) List of notification backend names.

## Attributes Reference

- **`id`** - The message type string (same as `message_type`).
- **`message_type_label`** - The human-readable label for the message type.
- **`receivers`** - The resolved list of receiver display names.
