---
weight: 6
title: "Add arguments to Terraform commands"
description: Tfautomv allows passing arbitrary arguments to Terraform init and plan commands.
---

# Add arguments to Terraform commands

Under the hood, `tfautomv` runs `terraform init` and `terraform plan`. In order
to pass addition arguments to these commands you should use Terraform's built-in
[`TF_CLI_ARGS` and `TF_CLI_ARGS_name` environment variables.](https://www.terraform.io/cli/config/environment-variables#tf_cli_args-and-tf_cli_args_name).

For example, in order to use a file of variables during Terraform's plan:

```bash
TF_CLI_ARGS_plan="-var-file=production.tfvars" tfautomv
```

You can also skip Terraform's refresh to speed up the planning step:

```bash
TF_CLI_ARGS_plan="-refresh=false" tfautomv
```
