---
weight: 7
title: "Use a specific Terraform binary"
description: Tfautomv allows using any Terraform binary or wrapper.
---

# Use a specific Terraform binary

Under the hood, `tfautomv` runs `terraform`. You can replace `terraform` with
any other binary with the `-terraform-bin` flag.

For example, you can use Terragrunt as a Terraform wrapper:

```bash
tfautomv -terraform-bin=terragrunt
```
