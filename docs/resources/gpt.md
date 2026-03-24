---
page_title: "jumpserver_gpt Resource - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Manages a GPT asset in JumpServer.
---

# `jumpserver_gpt` Resource

The `jumpserver_gpt` resource allows you to create and manage GPT-type assets in JumpServer (e.g. ChatGPT, LLM endpoints).

## Example Usage

```hcl
resource "jumpserver_gpt" "chatgpt" {
  name      = "chatgpt-prod"
  address   = "api.openai.com"
  platform  = 50
  zone_name = "Production"
  node_name = "AI Services"

  protocols {
    name = "https"
    port = 443
  }
}
```

## Argument Reference

* `name` - (Required) The name of the GPT asset.
* `address` - (Required) The address of the GPT endpoint.
* `platform` - (Required) The platform code for this GPT asset.
* `comment` - (Optional) A description for the GPT asset.
* `zone_name` - (Required) The name of the Zone this GPT belongs to.
* `node_name` - (Required) The name of the Node this GPT belongs to.
* `accounts` - (Optional) A list of account definitions. See the `jumpserver_host` resource for sub-attributes.
* `protocols` - (Optional) A list of protocols. See the `jumpserver_host` resource for sub-attributes.

## Attribute Reference

* `id` - The UUID of the GPT asset.
