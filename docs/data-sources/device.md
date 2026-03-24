---
page_title: "jumpserver_device Data Source - terraform-provider-jumpserver"
subcategory: "Assets"
description: |-
  Look up a device asset in JumpServer by name.
---

# `jumpserver_device` Data Source

Use this data source to look up an existing device asset in JumpServer by name.

## Example Usage

```hcl
data "jumpserver_device" "example" {
  name = "switch-01"
}
```

## Argument Reference

- **`name`** - (Required) The name of the device to look up.

## Attributes Reference

- **`id`** - The UUID of the device.
- **`address`** - The address of the device.
- **`comment`** - A comment or description.
- **`is_active`** - Whether the device is active.
