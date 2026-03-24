---
page_title: "jumpserver_session_join_record Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages a session join record in JumpServer.
---

# `jumpserver_session_join_record` Resource

The `jumpserver_session_join_record` resource allows you to manage session join records in JumpServer.

## Example Usage

```hcl
resource "jumpserver_session_join_record" "example" {
  session     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  sharing     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  verify_code = "123456"
}
```

## Argument Reference

- **`session`** - (Required) The session UUID.
- **`sharing`** - (Required) The session sharing UUID.
- **`verify_code`** - (Required) The verification code.
- **`joiner`** - (Optional) The joiner user UUID.
- **`remote_addr`** - (Optional) The remote address.
- **`login_from`** - (Optional) The login source: ST, RT, WT. Default: `WT`.
- **`reason`** - (Optional) The reason for joining.
- **`is_success`** - (Optional) Whether the join was successful. Default: `true`.
- **`is_finished`** - (Optional) Whether the join session is finished. Default: `false`.

## Attributes Reference

- **`id`** - The UUID of the join record.
- **`joiner_display`** - The joiner display name.
- **`action_permission`** - The action permission.
- **`created_by`** - The creator.
- **`date_joined`** - The join date.
- **`date_left`** - The leave date.
- **`date_created`** - The creation date.
