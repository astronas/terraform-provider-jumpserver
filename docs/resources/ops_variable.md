---
page_title: "jumpserver_ops_variable Resource - terraform-provider-jumpserver"
subcategory: "Ops"
description: |-
  Manages an ops variable in JumpServer.
---

# `jumpserver_ops_variable` Resource

The `jumpserver_ops_variable` resource allows you to create and manage variables used in JumpServer operations and playbooks.

## Example Usage

```hcl
resource "jumpserver_ops_variable" "env" {
  name          = "Environment"
  var_name      = "ENVIRONMENT"
  type          = "string"
  default_value = "production"
  tips          = "The deployment environment"
  required      = true
  comment       = "Environment variable"
}
```

## Argument Reference

* `name` - (Required) The display name of the variable.
* `var_name` - (Required) The variable name (used in scripts).
* `type` - (Required) The type of the variable.
* `default_value` - (Optional) The default value.
* `tips` - (Optional) Help text for the variable.
* `required` - (Optional) Whether the variable is required. Default: `false`.
* `comment` - (Optional) A description for the variable.

## Attribute Reference

* `id` - The UUID of the variable.
