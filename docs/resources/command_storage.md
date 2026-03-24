---
page_title: "jumpserver_command_storage Resource - terraform-provider-jumpserver"
subcategory: "Infrastructure"
description: |-
  Manages a JumpServer Command Storage.
---

# jumpserver_command_storage

Manages a JumpServer Command Storage backend. Defines where session command logs are stored (local server, Elasticsearch, etc.).

## Example Usage

```hcl
resource "jumpserver_command_storage" "elasticsearch" {
  name = "es-commands"
  type = "es"
  meta = jsonencode({
    HOSTS      = ["https://es.example.com:9200"]
    INDEX      = "jumpserver-commands"
    DOC_TYPE   = "_doc"
  })
  is_default = true
  comment    = "Elasticsearch command log storage"
}

resource "jumpserver_command_storage" "local" {
  name = "local-commands"
  type = "server"
  meta = jsonencode({})
}
```

## Argument Reference

- `name` - (Required) The name of the command storage.
- `type` - (Required, ForceNew) Storage type: `null`, `server`, `es` (Elasticsearch).
- `meta` - (Optional) Configuration as JSON string. Defaults to `{}`.
- `is_default` - (Optional) Whether this is the default command storage. Defaults to `false`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the command storage.

## Import

```shell
terraform import jumpserver_command_storage.example <uuid>
```
