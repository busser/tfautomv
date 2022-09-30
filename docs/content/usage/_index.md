---
weight: 2
title: "Usage"
---

# Usage

You can run `tfautomv` in any directory where you run `terraform plan`.

Tfautomv runs `terraform init` and `terraform plan` in the current directory and
analyses Terraform's plan. Based on this analysis, it identifies resources that
have moved in your codebase but not yet in Terraform's state. For each of those
resources, it appends a `moved` block to the `moves.tf` file.

Tfautomv is fully provider-agnostic. It works with all Terraform providers.

## How-to's

{{<section>}}

## Flags

```console
$ tfautomv -h
Usage of tfautomv:
  -dry-run
    	print moves instead of writing them to disk
  -ignore rule
    	ignore differences based on a rule
  -no-color
    	disable color in output
  -output format
    	output format of moves ("blocks" or "commands") (default "blocks")
  -show-analysis
    	show detailed analysis of Terraform plan
  -version
    	print version and exit
```
