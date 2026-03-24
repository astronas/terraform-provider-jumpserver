---
page_title: "jumpserver_applet_host_deployment Resource - terraform-provider-jumpserver"
subcategory: "Terminal"
description: |-
  Manages an applet host deployment in JumpServer.
---

# `jumpserver_applet_host_deployment` Resource

Manages applet host deployments in JumpServer. Triggers deployment of applets to a specified host.

## Example Usage

```hcl
resource "jumpserver_applet_host_deployment" "example" {
  host            = jumpserver_applet_host.example.id
  install_applets = true
  comment         = "Initial deployment"
}
```

## Argument Reference

- **`host`** - (Required) The UUID of the applet host to deploy to.
- **`install_applets`** - (Optional) Whether to install applets during deployment. Defaults to `true`.
- **`comment`** - (Optional) A description for the deployment.

## Attribute Reference

- **`id`** - The UUID of the deployment.
- **`date_created`** - The creation date.
- **`date_start`** - The deployment start date.
- **`date_finished`** - The deployment finish date.
