---
page_title: "jumpserver_chatai_prompt Resource - terraform-provider-jumpserver"
subcategory: "Settings"
description: |-
  Manages a ChatAI prompt in JumpServer.
---

# `jumpserver_chatai_prompt` Resource

The `jumpserver_chatai_prompt` resource allows you to create and manage ChatAI prompts in JumpServer.

## Example Usage

```hcl
resource "jumpserver_chatai_prompt" "example" {
  name    = "security-check"
  content = "Analyze the following command for security risks: {command}"
}
```

## Argument Reference

- **`name`** - (Required) The name of the ChatAI prompt.
- **`content`** - (Required) The content of the prompt.
- **`builtin`** - (Optional) Whether the prompt is builtin. Default: `false`.

## Attributes Reference

- **`id`** - The UUID of the prompt.
- **`date_created`** - The creation date.
- **`date_updated`** - The last update date.
