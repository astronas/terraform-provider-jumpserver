---
page_title: "jumpserver_job Resource - terraform-provider-jumpserver"
subcategory: "Operations"
description: |-
  Manages a job in JumpServer.
---

# jumpserver_job

Manages a job in JumpServer. Jobs allow executing ad-hoc commands, playbooks, or file uploads on target assets.

## Example Usage

### Ad-hoc Command

```hcl
resource "jumpserver_job" "check_disk" {
  name   = "check-disk-space"
  type   = "adhoc"
  module = "shell"
  args   = "df -h"
  runas  = "root"
  assets = ["asset-uuid-1", "asset-uuid-2"]
}
```

### Playbook Job

```hcl
resource "jumpserver_job" "deploy" {
  name     = "deploy-app"
  type     = "playbook"
  playbook = jumpserver_playbook.deploy.id
  assets   = ["asset-uuid-1"]

  is_periodic = true
  crontab     = "0 2 * * *"
  comment     = "Nightly deployment"
}
```

## Argument Reference

* `name` - (Required) The name of the job.
* `type` - (Optional) The type of job: `adhoc`, `playbook`, or `upload_file`. Default: `adhoc`. Forces new resource.
* `module` - (Optional) The module to use for adhoc jobs: `shell`, `winshell`, `python`, `raw`. Default: `shell`.
* `args` - (Optional) Command or arguments to execute.
* `playbook` - (Optional) The UUID of the playbook to run (for playbook type jobs).
* `assets` - (Optional) List of asset UUIDs to run the job on.
* `nodes` - (Optional) List of node UUIDs to run the job on.
* `runas_policy` - (Optional) Run-as policy: `skip` or `if_need`. Default: `skip`.
* `runas` - (Optional) The user to run as on the target assets. Default: `root`.
* `timeout` - (Optional) Timeout in seconds (-1 for no timeout). Default: `-1`.
* `chdir` - (Optional) Working directory on the target.
* `is_periodic` - (Optional) Whether the job runs periodically. Default: `false`.
* `interval` - (Optional) Interval in hours for periodic execution (1-65535). Default: `24`.
* `crontab` - (Optional) Crontab expression for scheduling.
* `run_after_save` - (Optional) Whether to run the job immediately after saving. Default: `false`.
* `use_parameter_define` - (Optional) Whether to use parameter definitions. Default: `false`.
* `parameters_define` - (Optional) JSON string of parameter definitions.
* `comment` - (Optional) A description for the job.

## Attribute Reference

* `id` - The UUID of the job.
