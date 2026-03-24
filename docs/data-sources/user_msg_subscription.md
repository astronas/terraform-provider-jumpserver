---
page_title: "jumpserver_user_msg_subscription Data Source - terraform-provider-jumpserver"
subcategory: "Notifications"
description: |-
  Look up a user's message subscription in JumpServer.
---

# `jumpserver_user_msg_subscription` Data Source

Use this data source to look up the notification backends a user is subscribed to in JumpServer.

## Example Usage

```hcl
data "jumpserver_user_msg_subscription" "admin" {
  user_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

* `user_id` - (Required) The user UUID to look up subscriptions for.

## Attribute Reference

* `id` - The user UUID.
* `receive_backends` - List of notification backend names the user is subscribed to.
