---
page_title: "jumpserver_setting Resource - terraform-provider-jumpserver"
subcategory: "Settings"
description: |-
  Manages basic settings in JumpServer.
---

# `jumpserver_setting` Resource

The `jumpserver_setting` resource allows you to manage JumpServer basic settings. This is a singleton resource — only one instance can exist.

## Example Usage

```hcl
resource "jumpserver_setting" "basic" {
  site_url               = "https://jumpserver.example.com"
  user_guide_url         = "https://docs.example.com/guide"
  global_org_display_name = "MyOrganization"
}
```

## Argument Reference

- **`site_url`** - (Required) The site URL (URI format).
- **`user_guide_url`** - (Optional) The user guide URL.
- **`global_org_display_name`** - (Optional) The global organization display name.
- **`help_document_url`** - (Optional) The help document URL.
- **`help_support_url`** - (Optional) The help support URL.

## Attributes Reference

- **`id`** - Always `settings` (singleton resource).
