---
page_title: "jumpserver_endpoint_rule Resource - terraform-provider-jumpserver"
subcategory: "Infrastructure"
description: |-
  Manages a JumpServer Endpoint Rule.
---

# jumpserver_endpoint_rule

Manages a JumpServer Endpoint Rule. Rules route traffic from specific IP ranges to designated endpoints with priority-based matching.

## Example Usage

```hcl
resource "jumpserver_endpoint_rule" "internal" {
  name     = "internal-routing"
  ip_group = ["10.0.0.0/8", "192.168.0.0/16"]
  priority = 10
  endpoint = jumpserver_endpoint.internal.id
  comment  = "Route internal traffic to internal endpoint"
}

resource "jumpserver_endpoint_rule" "default" {
  name     = "default-routing"
  ip_group = ["*"]
  priority = 100
  endpoint = jumpserver_endpoint.main.id
  comment  = "Default route for all other traffic"
}
```

## Argument Reference

- `name` - (Required) The name of the endpoint rule.
- `ip_group` - (Optional) IP ranges to match. Supports CIDR, ranges, and wildcards. Defaults to `["*"]`.
- `priority` - (Optional) Priority 1-100. Lower values match first. Defaults to `50`.
- `endpoint` - (Optional) Endpoint UUID to route matching traffic to.
- `is_active` - (Optional) Whether the rule is active. Defaults to `true`.
- `comment` - (Optional) A description.

## Attribute Reference

- `id` - The UUID of the endpoint rule.

## Import

```shell
terraform import jumpserver_endpoint_rule.example <uuid>
```
