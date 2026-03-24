---
page_title: "jumpserver_ssh_key Resource - terraform-provider-jumpserver"
subcategory: "Authentication"
description: |-
  Manages an SSH key in JumpServer.
---

# `jumpserver_ssh_key` Resource

Manages SSH keys for user authentication in JumpServer.

## Example Usage

```hcl
resource "jumpserver_ssh_key" "example" {
  name       = "my-ssh-key"
  public_key = file("~/.ssh/id_ed25519.pub")
  is_active  = true
  comment    = "My workstation key"
}
```

## Argument Reference

- **`name`** - (Required) The name of the SSH key.
- **`public_key`** - (Optional) The public key content.
- **`is_active`** - (Optional) Whether the key is active. Defaults to `true`.
- **`comment`** - (Optional) A comment or description.

## Attribute Reference

- **`id`** - The UUID of the SSH key.
- **`public_key_comment`** - The comment embedded in the public key.
- **`public_key_hash_md5`** - The MD5 hash of the public key.
- **`date_last_used`** - The date the key was last used.
- **`date_created`** - The creation date.
