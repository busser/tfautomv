# tfautomv <!-- omit in toc -->

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub release](https://img.shields.io/github/release/busser/tfautomv.svg)](https://github.com/busser/tfautomv/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/busser/tfautomv)](https://goreportcard.com/report/github.com/busser/tfautomv)

Generate Terraform `moved` blocks automatically.

- [Why?](#why)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Homebrew](#homebrew)
  - [Yay](#yay)
  - [Shell script](#shell-script)
  - [Download](#download)
  - [From source](#from-source)
- [Usage](#usage)
  - [Generating `moved` blocks](#generating-moved-blocks)
  - [Generating `terraform state mv` commands](#generating-terraform-state-mv-commands)
  - [Understanding why a resource was not matched](#understanding-why-a-resource-was-not-matched)
  - [Ignoring certain differences](#ignoring-certain-differences)
    - [The `everything` kind](#the-everything-kind)
    - [The `whitespace` kind](#the-whitespace-kind)
    - [The `prefix` kind](#the-prefix-kind)
    - [Referencing nested attributes](#referencing-nested-attributes)
  - [Passing additional arguments to Terraform](#passing-additional-arguments-to-terraform)
  - [Using Terragrunt instead of Terraform](#using-terragrunt-instead-of-terraform)
  - [Disabling colors in output](#disabling-colors-in-output)
- [Thanks](#thanks)
- [License](#license)

## Why?

`tfautomv` (a.k.a Terraform auto-move) is a refactoring helper. With it, making
structural changes to your Terraform codebase becomes much easier.

When you move a resource in your code, Terraform loses track of the resource's
state. The next time you run Terraform, it will plan to delete the resource it
has memory of and create the "new" resource it found in your refactored code.

`tfautomv` inspects the output of `terraform plan`, detects such
creation/deletion pairs and writes a [`moved` block](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring#moved-block-syntax)
so that Terraform now knows no deletion or creation is required.

We explain why we built tfautomv in more detail [in this blog article](https://www.padok.fr/en/blog/terraform-refactoring-tfautomv).

Here's a quick view of what `tfautomv` does:

![demo](./docs/content/getting-started/demo.gif)

## Requirements

`tfautomv` uses the Terraform CLI command under the hood. This allows it to work
with any Terraform version reliably.

You will need Terraform v1.1 or above to generate `moved` blocks. Or Terraform
v0.13 or above to generate `terraform state mv` commands.

## Installation

_Contributions to support other installation methods are welcome!_

### Homebrew

On MacOS or Linux:

```bash
brew install busser/tap/tfautomv
```

### Yay

On Arch Linux:

```bash
yay tfautomv-bin
```

### Shell script

On MacOS or Linux:

```bash
curl -sSfL https://raw.githubusercontent.com/busser/tfautomv/main/install.sh | sh
```

_This script can probably support Windows with a small amount of work.
Contributions welcome!_

### Download

On the Github repository's [Releases page](https://github.com/busser/tfautomv/releases),
download the binary that matches your workstation's OS and CPU architecture.

Put the binary in a directory present in your system's `PATH` environment
variable.

### From source

You must have Go 1.18+ installed to compile tfautomv.

Clone the repository and build the binary:

```bash
git clone https://github.com/busser/tfautomv
cd tfautomv
make build
```

Then, move `bin/tfautomv` to a directory resent in your system's `PATH`
environment variable.

## Usage

### Generating `moved` blocks

In any directory where you would run `terraform plan`, you can run:

```bash
tfautomv
```

This will run `terraform init`, `terraform plan`, and then write `moved` blocks
to a `moves.tf` file.

That's all there is to it!

### Generating `terraform state mv` commands

If you are using a version of Terraform older than v1.1 or don't want to use
`moved` blocks, you can generate `terraform state mv` commands instead:

```bash
tfautomv -output=commands
```

This will print commands to standard output. You can copy and paste them to a
terminal to run them manually.

Alternatively, you can write the commands to a file:

```bash
tfautomv -output=commands > moves.sh
```

Or pipe them into a shell to run them immediately:

```bash
tfautomv -output=commands | sh
```

### Understanding why a resource was not matched

If you are not seeing `moved` blocks for a resource you expected to be matched,
you can run `tfautomv` with the `-show-analysis` flag to get more information:

```bash
tfautomv -show-analysis
```

This will print a detailed analysis of why each resource was or was not matched.

From there, you can choose to edit your code, write a `moved` block manually, or
use the `-ignore` flag to ignore certain differences.

### Ignoring certain differences

`tfautomv` works by comparing resources Terraform plans to create (those in your
code) to those Terraform plans to delete (those in your state). Sometimes,
`tfautomv` may not be able to match two resources together because of a
difference in a specific attribute, even though the resources are in fact the
same. This usually happens when the Terraform provider that manages the resource
has transformed the attribute's value in some way.

In those cases, you can use the `-ignore` flag to ignore specific differences.
`tfautomv` will ignore differences based on a set of rules that you can
provide.

Each rule includes:

- A _kind_ that identifies the nature of the difference to ignore
- A _resource type_ the rule applies to
- An _attribute_ inside the resource the rule applies to
- Optionally, additional arguments specific to the class

A rule is written as a colon-separated string:

```plaintext
<KIND>:<RESOURCE TYPE>:<ATTRIBUTE NAME>[:<KIND ARGUMENTS>]
```

You can use the `-ignore` flag multiple times to provide multiple rules:

```bash
tfautomv \
  -ignore="whitespace:azurerm_api_management_policy:xml_content" \
  -ignore="prefix:google_storage_bucket_iam_member:bucket:b/"
```

_If you have a use case that is not covered by existing kinds, please open an
issue so we can track demand for it._

#### The `everything` kind

Use the `everything` kind to ignore any difference between two values of an
attribute:

```bash
tfautomv -ignore="everything:<RESOURCE TYPE>:<ATTRIBUTE>"
```

For example:

```bash
tfautomv -ignore="everything:random_pet:length"
```

#### The `whitespace` kind

Use the `whitespace` kind to ignore differences in whitespace between two
values of an attribute:

```bash
tfautomv -ignore="whitespace:<RESOURCE TYPE>:<ATTRIBUTE NAME>"
```

For example, this rule:

```bash
tfautomv -ignore="whitespace:azurerm_api_management_policy:xml_content"
```

will allow these two resources to match:

```terraform
# This resource has its XML nicely formatted.
resource "azurerm_api_management_policy" "foo" {
  api_management_id = "..."

  xml_content = <<-EOT
  <policies>
    <inbound>
      <cross-domain />
      <base />
      <find-and-replace from="xyz" to="abc" />
    </inbound>
  </policies>
  EOT
}

# This resource has its XML on one line.
resource "azurerm_api_management_policy" "bar" {
  api_management_id = "..."

  xml_content = "<policies><inbound><cross-domain /><base /><find-and-replace from=\"xyz\" to=\"abc\" /></inbound></policies>"
}
```

#### The `prefix` kind

Use the `prefix` kind to ignore a specific prefix between in one of two values
of an attribute:

```bash
tfautomv -ignore="prefix:<RESOURCE TYPE>:<ATTRIBUTE NAME>:<PREFIX>"
```

For example:

```bash
tfautomv -ignore="prefix:google_storage_bucket_iam_member:bucket:b/"
```

will strip the `b/` prefix from the `bucket` attribute of any
`google_storage_bucket_iam_member` resources before comparing the attirbute's
values.

#### Referencing nested attributes

Join parent attributes with child attributes with a `.`:

```plaintext
<KIND>:<RESOURCE TYPE>:parent_obj.child_field
<KIND>:<RESOURCE TYPE>:parent_list.0
```

If using the `-show-analysis` flag, you can see the full path to an attribute in
the analysis output.

### Passing additional arguments to Terraform

You can pass additional arguments to Terraform by using Terraform's built-in
[`TF_CLI_ARGS` and `TF_CLI_ARGS_name` environment variables.](https://www.terraform.io/cli/config/environment-variables#tf_cli_args-and-tf_cli_args_name).

For example, in order to use a file of variables during Terraform's plan:

```bash
TF_CLI_ARGS_plan="-var-file=production.tfvars" tfautomv
```

Or to skip Terraform's refresh and speed up the planning step:

```bash
TF_CLI_ARGS_plan="-refresh=false" tfautomv
```

### Using Terragrunt instead of Terraform

You can tell `tfautomv` to use the Terragrunt CLI instead of the Terraform CLI
with the `-terraform-bin` flag:

```bash
tfautomv -terraform-bin=terragrunt
```

This also works with any other executable that has an `init` and `plan` command.

### Disabling colors in output

Add the `-no-color` flag to your `tfautomv` command to disable output
formatting like colors, bold text, etc.

For example:

```bash
tfautomv -no-color
```

## Thanks

Thanks to [Padok](https://www.padok.fr), where this project was born 💜

## License

The code is licensed under the permissive Apache v2.0 license. [Read this](<https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)>) for a summary.
