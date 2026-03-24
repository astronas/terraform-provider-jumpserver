---
page_title: "jumpserver_system_user Resource - terraform-provider-jumpserver"
subcategory: "Accounts"
description: |-
  Manages a system user in JumpServer.
---

# `jumpserver_system_user` Resource

The jumpserver_system_user resource allows you to create and manage system users in Jumpserver. A system user is a user account that is used to access assets

## Example Usage

```hcl
resource "jumpserver_system_user" "example_user" {
  name                    = "student"
  username                = "student"
  password                = "studentpass"
  type                    = "common"
  protocol                = "ssh"
  login_mode              = "auto"
  priority                = 81
  sudo                    = "/bin/whoami"
  shell                   = "/bin/bash"
  sftp_root               = "tmp"
  home                    = "/home/student"
  username_same_with_user = false
  auto_push               = false
  su_enabled              = false
}
```

## Argument Reference

* `name` - (Required) The name of the system user.
* `username` - (Optional) The username of the system user.
* `password` - (Optional) The password of the system user.
* `type` - (Optional) The type of system user (e.g., common).
* `protocol` - (Optional) The protocol the system user will use.
* `login_mode` - (Optional) The login mode of the system user.
* `priority` - (Optional) The priority of the system user.
* `sudo` - (Optional) The sudo command for the system user.
* `shell` - (Optional) The shell for the system user.
* `sftp_root` - (Optional) The SFTP root directory for the system user.
* `home` - (Optional) The home directory for the system user.
* `username_same_with_user` - (Optional) Whether the username is the same as the user.
* `auto_push` - (Optional) Whether to auto-push the system user.
* `su_enabled` - (Optional) Whether the system user can use su.

## Attribute Reference

* `id` - The ID of the system user.
* `name` - The name of the system user.
* `username` - The username of the system user.
* `type` - The type of system user.
* `protocol` - The protocol the system user will use.
* `login_mode` - The login mode of the system user.
* `priority` - The priority of the system user.
* `sudo` - The sudo command for the system user.
* `shell` - The shell for the system user.
* `sftp_root` - The SFTP root directory for the system user.
* `home` - The home directory for the system user.
* `username_same_with_user` - Whether the username is the same as the user.
* `auto_push` - Whether to auto-push the system user.
* `su_enabled` - Whether the system user can use su.