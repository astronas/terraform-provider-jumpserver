---
page_title: "jumpserver_playbook Resource - terraform-provider-jumpserver"
subcategory: "Operations"
description: |-
  Manages a playbook in JumpServer.
---

# jumpserver_playbook

Manages a playbook in JumpServer. Playbooks can be created blank or imported from a VCS repository, and used in jobs for automation.

## Example Usage

### Blank Playbook

```hcl
resource "jumpserver_playbook" "deploy" {
  name          = "deploy-application"
  scope         = "public"
  create_method = "blank"
  comment       = "Deployment playbook"
}
```

### VCS Playbook

```hcl
resource "jumpserver_playbook" "from_git" {
  name          = "infra-setup"
  scope         = "private"
  create_method = "vcs"
  vcs_url       = "https://git.example.com/playbooks/infra.git"
  comment       = "Infrastructure setup from VCS"
}
```

## Argument Reference

* `name` - (Required) The name of the playbook.
* `scope` - (Optional) The scope of the playbook: `public` or `private`. Default: `public`.
* `create_method` - (Optional) The creation method: `blank` or `vcs`. Default: `blank`. Forces new resource.
* `vcs_url` - (Optional) VCS repository URL (used when create_method is `vcs`).
* `comment` - (Optional) A description for the playbook.

## Attribute Reference

* `id` - The UUID of the playbook.
* `path` - The path to the playbook on the server.
