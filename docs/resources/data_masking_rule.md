---
page_title: "jumpserver_data_masking_rule Resource - terraform-provider-jumpserver"
subcategory: "Permissions & ACLs"
description: |-
  Manages a JumpServer Data Masking Rule.
---

# jumpserver_data_masking_rule

Manages a JumpServer Data Masking Rule. Masks sensitive data fields in session outputs.

## Example Usage

```hcl
resource "jumpserver_data_masking_rule" "mask_passwords" {
  name           = "mask-password-fields"
  users          = jsonencode({ type = "all" })
  assets         = jsonencode({ type = "all" })
  accounts       = ["@ALL"]
  priority       = 10
  action         = "accept"
  fields_pattern = "password|secret|token"
  masking_method = "fixed_char"
  mask_pattern   = "######"
  is_active      = true
  comment        = "Mask password-related fields"
}
```

## Argument Reference

- `name` - (Required) The name of the data masking rule.
- `users` - (Required) User filter as JSON string.
- `assets` - (Required) Asset filter as JSON string.
- `accounts` - (Required) List of account usernames.
- `priority` - (Optional) Priority 1-100. Defaults to `50`.
- `action` - (Optional) Action: `reject`, `accept`, `review`. Defaults to `reject`.
- `reviewers` - (Optional) List of reviewer user UUIDs.
- `fields_pattern` - (Optional) Regex pattern matching fields to mask. Defaults to `password`.
- `masking_method` - (Optional) Masking method: `fixed_char`. Defaults to `fixed_char`.
- `mask_pattern` - (Optional) Pattern used to mask fields. Defaults to `######`.
- `is_active` - (Optional) Whether the rule is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the data masking rule.

## Import

```shell
terraform import jumpserver_data_masking_rule.example <uuid>
```
