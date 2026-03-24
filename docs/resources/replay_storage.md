---
page_title: "jumpserver_replay_storage Resource - terraform-provider-jumpserver"
subcategory: "Infrastructure"
description: |-
  Manages a JumpServer Replay Storage.
---

# jumpserver_replay_storage

Manages a JumpServer Replay Storage backend. Defines where session replay recordings are stored (S3, Azure Blob, SFTP, etc.).

## Example Usage

```hcl
resource "jumpserver_replay_storage" "s3" {
  name = "s3-replays"
  type = "s3"
  meta = jsonencode({
    BUCKET     = "jumpserver-replays"
    ACCESS_KEY = var.aws_access_key
    SECRET_KEY = var.aws_secret_key
    REGION     = "us-east-1"
  })
  is_default = true
  comment    = "S3 replay storage"
}

resource "jumpserver_replay_storage" "azure" {
  name = "azure-replays"
  type = "azure"
  meta = jsonencode({
    CONTAINER_NAME = "replays"
    ACCOUNT_NAME   = "jumpserverstorage"
    ACCOUNT_KEY    = var.azure_storage_key
  })
}

resource "jumpserver_replay_storage" "sftp" {
  name = "sftp-replays"
  type = "sftp"
  meta = jsonencode({
    HOST     = "sftp.example.com"
    PORT     = 22
    USERNAME = "jumpserver"
    PASSWORD = var.sftp_password
    DIR      = "/replays"
  })
}
```

## Argument Reference

- `name` - (Required) The name of the replay storage.
- `type` - (Required, ForceNew) Storage type: `null`, `server`, `s3`, `ceph`, `swift`, `oss`, `azure`, `obs`, `cos`, `sftp`.
- `meta` - (Optional) Configuration as JSON string. Defaults to `{}`.
- `is_default` - (Optional) Whether this is the default replay storage. Defaults to `false`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the replay storage.

## Import

```shell
terraform import jumpserver_replay_storage.example <uuid>
```
