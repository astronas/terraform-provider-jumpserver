---
page_title: "jumpserver_cloud Resource"
subcategory: "Assets"
description: |-
  Manages a cloud asset in JumpServer.
---

# jumpserver_cloud

Manages a cloud asset in JumpServer. Used for cloud service endpoints such as AWS, Azure, Google Cloud APIs.

## Example Usage

```hcl
resource "jumpserver_cloud" "aws_account" {
  name      = "aws-prod"
  address   = "https://console.aws.amazon.com"
  platform  = 20
  zone_name = "Cloud"
  node_name = "AWS"
  comment   = "AWS production account"

  accounts {
    name        = "aws-admin"
    username    = "admin"
    secret_type = "access_key"
    secret      = "AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}
```

## Argument Reference

- `name` - (Required) The name of the cloud asset.
- `address` - (Required) The URL of the cloud service endpoint.
- `platform` - (Required) The platform ID (integer).
- `zone_name` - (Required) The name of the zone.
- `node_name` - (Required) The name of the node.
- `comment` - (Optional) A description.
- `is_active` - (Optional) Whether the asset is active. Default: `true`.
- `accounts` - (Optional) List of account blocks.
- `protocols` - (Optional) List of protocol blocks (`name`, `port`).

## Attribute Reference

- `id` - The UUID of the cloud asset.
- `zone_id` - The resolved zone UUID.
- `node_ids` - The resolved node UUIDs.
