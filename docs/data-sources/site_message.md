---
page_title: "jumpserver_site_message Data Source - terraform-provider-jumpserver"
subcategory: "Notifications"
description: |-
  Look up a site message in JumpServer by subject.
---

# `jumpserver_site_message` Data Source

Use this data source to look up a site message in JumpServer by its subject.

## Example Usage

```hcl
data "jumpserver_site_message" "latest" {
  subject = "Account change notification"
}
```

## Argument Reference

* `subject` - (Required) The subject of the site message to look up.

## Attribute Reference

* `id` - The UUID of the site message.
* `has_read` - Whether the message has been read.
* `read_at` - The date the message was read.
* `content` - The message content (JSON).
* `date_created` - The creation date.
