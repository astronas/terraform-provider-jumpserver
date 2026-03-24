---
page_title: "jumpserver_sms_backend Data Source - terraform-provider-jumpserver"
subcategory: "Settings"
description: |-
  Look up an SMS backend in JumpServer by name.
---

# `jumpserver_sms_backend` Data Source

Use this data source to look up an SMS backend configuration in JumpServer by its name.

## Example Usage

```hcl
data "jumpserver_sms_backend" "aliyun" {
  name = "alibaba_cloud"
}
```

## Argument Reference

* `name` - (Required) The name of the SMS backend to look up.

## Attribute Reference

* `id` - The name of the SMS backend.
* `label` - The human-readable label of the backend.
