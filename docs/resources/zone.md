---
page_title: "jumpserver_zone Resource"
subcategory: "Assets"
description: |-
  Manages a zone (domain) in JumpServer.
---

# jumpserver_zone

Manages a zone in JumpServer. Zones are used to group assets by network domain and associate gateways for connectivity.

## Example Usage

```hcl
resource "jumpserver_zone" "dmz" {
  name    = "DMZ"
  comment = "Demilitarized zone for public-facing servers"
}
```

## Argument Reference

- `name` - (Required) The name of the zone.
- `comment` - (Optional) A description for the zone.

## Attribute Reference

- `id` - The UUID of the zone.
