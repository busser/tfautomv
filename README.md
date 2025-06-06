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
  - [asdf](#asdf)
  - [Shell script](#shell-script)
  - [Download](#download)
  - [From source](#from-source)
- [Usage](#usage)
  - [Generating `moved` blocks](#generating-moved-blocks)
  - [Generating `terraform state mv` commands](#generating-terraform-state-mv-commands)
  - [Finding moves across multiple directories](#finding-moves-across-multiple-directories)
  - [Skipping the `init` and `refresh` steps](#skipping-the-init-and-refresh-steps)
  - [Understanding why a resource was not moved](#understanding-why-a-resource-was-not-moved)
  - [Ignoring certain differences](#ignoring-certain-differences)
    - [The `everything` kind](#the-everything-kind)
    - [The `whitespace` kind](#the-whitespace-kind)
    - [The `prefix` kind](#the-prefix-kind)
    - [Referencing nested attributes](#referencing-nested-attributes)
  - [Passing additional arguments to Terraform](#passing-additional-arguments-to-terraform)
  - [Using Terragrunt instead of Terraform](#using-terragrunt-instead-of-terraform)
  - [Using OpenTofu instead of Terraform](#using-opentofu-instead-of-terraform)
  - [Using other Terraform-compatible tools](#using-other-terraform-compatible-tools)
  - [Using existing plan files](#using-existing-plan-files)
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

## Best Practices

`tfautomv` is designed for **refactoring scenarios** where you want to restructure your Terraform code without changing the actual infrastructure. Understanding this distinction is crucial for successful usage.

### ✅ Good use cases (pure refactoring)

- **Renaming resources**: `aws_instance.web` → `aws_instance.web_server`
- **Moving resources between modules**: `aws_instance.web` → `module.ec2.aws_instance.web`
- **Changing resource organization**: Converting single resources to `for_each` loops
- **Module restructuring**: Moving resources between different modules

### ❌ Problematic use cases (infrastructure changes)

- **Changing resource configuration**: Removing or adding attributes like `tags`, `security_groups`, etc.
- **Combining refactoring with configuration changes**: Renaming AND modifying attributes in the same operation
- **Provider-specific transformations**: Changes where the provider modifies attribute values

### Recommended workflow

1. **First, make structural changes only**: Rename resources, move between modules, or restructure without changing any resource attributes
2. **Run tfautomv**: Generate the appropriate `moved` blocks or `terraform state mv` commands
3. **Apply the moves**: This should result in an empty or minimal plan showing no infrastructure changes
4. **Then make configuration changes**: In a separate step, modify resource attributes as needed

This approach ensures that moves represent safe refactoring operations separate from infrastructure changes.

## Requirements

`tfautomv` uses the Terraform CLI command under the hood. This allows it to work
with any Terraform version reliably.

Certain features require specific versions of Terraform:

- `moved` blocks require Terraform v1.1 or above
- cross-module `terraform state mv` commands require Terraform v0.14 or above
- single-module `terraform state mv` commands require Terraform v0.13 or above

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

### asdf

With asdf version manager:

```bash
asdf plugin add tfautomv https://github.com/busser/asdf-tfautomv.git
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

This will run `terraform init`, `terraform refresh`, and `terraform plan`, and
then write `moved` blocks to a `moves.tf` file.

That's all there is to it!

You can also target a specific working directory:

```bash
tfautomv ./production
```

### Generating `terraform state mv` commands

By default, tfautomv writes moves to as `moved` blocks when possible and falls
back to `terraform state mv` commands when not.

You can force `tfautomv` to write only `moved` blocks with the `--output=moves`
flag:

```bash
tfautomv --output=moves
```

You can force `tfautomv` to write only `terraform state mv` commands with the
`--output=commands` flag:

```bash
tfautomv --output=commands
```

This will print commands to standard output. You can copy and paste them to a
terminal to run them manually.

Alternatively, you can write the commands to a file:

```bash
tfautomv --output=commands > moves.sh
```

Or pipe them into a shell to run them immediately:

```bash
tfautomv --output=commands | sh
```

The `-o` flag is shorthand for `--output`:

```bash
tfautomv -o commands
```

### Finding moves across multiple directories

If you have multiple Terraform modules in different directories, you can pass
those directories to `tfautomv`:

```bash
tfautomv ./production/main ./production/backup -o commands
```

This will run `terraform init`, `terraform refresh`, and `terraform plan` in
each directory, and then write `terraform state mv` commands to standard output.
These commands will move resources within and across directories as needed.

Terraform does not natively support moving resources across directories. To
achieve this, `tfautomv` will output commands that pull copies of each
directory's state, perform the moves, and then push the new state back to the
directory's state backend.

You can pass as many directories as you want to `tfautomv`.

This is only compatible with the `commands` output format. Terraform's `moved`
block syntax does not support moving resources across directories.

### Skipping the `init` and `refresh` steps

By default, `tfautomv` runs Terraform's `init`, `refresh`, and `plan` steps.
To save time, you can skip the `init` or `refresh` steps with the `--skip-init`
and `--skip-refresh` flags:

```bash
tfautomv --skip-init --skip-refresh
```

The `-s` flag is shorthand for `--skip-init` and `-S` for `--skip-refresh:

```bash
tfautomv -sS
```

### Understanding why a resource was not moved

If you are not seeing a `moved` block for a resource you expected to be moved,
you can increase `tfautomv`'s verbosity with the `-v` flag to get more
information:

```bash
tfautomv -v
```

The default verbosity level is 0. You can increase the verbosity up to 3 by
repeating the `-v` flag:

```bash
tfautomv -vvv
```

Alternatively, you can specify a specific verbosity level with the `--verbosity`
flag:

```bash
tfautomv --verbosity=2
```

Based on why the resource was not moved, you can choose to edit your code,
write a `moved` block manually, or use the `-ignore` flag to ignore certain
differences.

|                     level 0 (default)                     |                      level 1 (`-v`)                       |                      level 2 (`-vv`)                      |                     level 3 (`-vvv`)                      |
| :-------------------------------------------------------: | :-------------------------------------------------------: | :-------------------------------------------------------: | :-------------------------------------------------------: |
| ![verbosity level 0](./docs/images/verbosity/level-0.png) | ![verbosity level 1](./docs/images/verbosity/level-1.png) | ![verbosity level 2](./docs/images/verbosity/level-2.png) | ![verbosity level 3](./docs/images/verbosity/level-3.png) |

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

#### ⚠️ Important considerations when using `--ignore`

The `--ignore` flag tells tfautomv to act as if certain attributes don't exist when comparing resources. While powerful, this comes with risks:

- **Risk of incorrect matches**: If ignored attributes are actually important for identifying the correct resource pairing, tfautomv may match unrelated resources or fail to find matches
- **Intended for provider quirks**: Use primarily when providers transform attribute values in ways that don't reflect actual infrastructure changes
- **Not for configuration changes**: Avoid using `--ignore` to force matches when you've intentionally changed resource configuration

**Good use cases for `--ignore`**:
- Provider transforms whitespace in policy documents
- Provider adds computed fields that weren't in the original configuration
- Provider normalizes values (e.g., adding default ports to security group rules)

**Problematic use cases for `--ignore`**:
- Forcing matches when you've intentionally changed tags, security groups, or other meaningful attributes
- Ignoring differences that represent real infrastructure changes you made

Each rule includes:

- A _kind_ that identifies the nature of the difference to ignore
- A _resource type_ the rule applies to
- An _attribute_ inside the resource the rule applies to
- Optionally, additional arguments specific to the class

A rule is written as a colon-separated string:

```plaintext
<KIND>:<RESOURCE TYPE>:<ATTRIBUTE NAME>[:<KIND ARGUMENTS>]
```

You can use the `--ignore` flag multiple times to provide multiple rules:

```bash
tfautomv \
  --ignore="whitespace:azurerm_api_management_policy:xml_content" \
  --ignore="prefix:google_storage_bucket_iam_member:bucket:b/"
```

_If you have a use case that is not covered by existing kinds, please open an
issue so we can track demand for it._

#### The `everything` kind

Use the `everything` kind to ignore any difference between two values of an
attribute:

```bash
tfautomv --ignore="everything:<RESOURCE TYPE>:<ATTRIBUTE>"
```

For example:

```bash
tfautomv --ignore="everything:random_pet:length"
```

#### The `whitespace` kind

Use the `whitespace` kind to ignore differences in whitespace between two
values of an attribute:

```bash
tfautomv --ignore="whitespace:<RESOURCE TYPE>:<ATTRIBUTE NAME>"
```

For example, this rule:

```bash
tfautomv --ignore="whitespace:azurerm_api_management_policy:xml_content"
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
tfautomv --ignore="prefix:<RESOURCE TYPE>:<ATTRIBUTE NAME>:<PREFIX>"
```

For example:

```bash
tfautomv --ignore="prefix:google_storage_bucket_iam_member:bucket:b/"
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

To get an attribute's full path, increase the verbosity level with the `-v`
flag:

```bash
tfautomv -vvv
```

#### Examples: When to use and avoid `--ignore`

**✅ Good example - Provider-transformed attribute**:
```bash
# Provider normalizes JSON policy formatting
tfautomv --ignore="whitespace:aws_iam_policy:policy"
```

**✅ Good example - Provider adds computed fields**:
```bash
# Provider adds computed "arn" or "id" fields that weren't in configuration
tfautomv --ignore="everything:aws_s3_bucket:arn"
```

**❌ Problematic example - Intentional configuration change**:
```bash
# DON'T do this - you've intentionally changed tags
# This forces a match that may result in unintended infrastructure changes
tfautomv --ignore="everything:aws_instance:tags"

# Instead: Apply moves first, then change tags in a separate operation
```

**❌ Problematic example - Mixing refactoring with changes**:
```terraform
# Before (applied to infrastructure):
resource "aws_instance" "web" {
  ami           = "ami-12345"
  instance_type = "t2.micro"
  tags = {
    Environment = "prod"
    Team        = "backend"
  }
}

# After (refactored AND changed):
resource "aws_instance" "web_server" {  # renamed
  ami           = "ami-12345"
  instance_type = "t2.micro"
  tags = {
    Environment = "production"  # changed value
    Team        = "backend"
    Project     = "website"     # added new tag
  }
}

# DON'T use --ignore here - this mixes refactoring with real changes
# Instead: rename first, apply moves, then modify tags separately
```

### Passing additional arguments to Terraform

You can pass additional arguments to Terraform by using Terraform's built-in
[`TF_CLI_ARGS` and `TF_CLI_ARGS_name` environment variables.](https://www.terraform.io/cli/config/environment-variables#tf_cli_args-and-tf_cli_args_name).

For example, in order to use a file of variables during Terraform's plan:

```bash
TF_CLI_ARGS_plan="-var-file=production.tfvars" tfautomv
```

### Using Terragrunt instead of Terraform

You can tell `tfautomv` to use the Terragrunt CLI instead of the Terraform CLI
with the `--terraform-bin` flag:

```bash
tfautomv --terraform-bin=terragrunt
```

### Using OpenTofu instead of Terraform

OpenTofu is officially supported! You can use OpenTofu with the `--terraform-bin` flag:

```bash
tfautomv --terraform-bin=tofu
```

This works with all tfautomv features including `moved` blocks, `tofu state mv` commands, and the `--preplanned` flag.

### Using other Terraform-compatible tools

The `--terraform-bin` flag works with any executable that has an `init` and `plan` command compatible with Terraform.

### Using existing plan files

If you have already generated Terraform plan files, you can use them directly with the `--preplanned` flag instead of having tfautomv run `terraform plan`. This is useful for:

- **Performance**: Avoid re-running expensive plan operations when iterating on `--ignore` rules
- **Enterprise environments**: Where running terraform locally is complex due to secrets or remote state
- **CI/CD workflows**: Where plans are generated in earlier pipeline stages
- **Remote workspaces**: TFE/Cloud environments where you can download JSON plans but can't run terraform locally

#### Single directory with default plan file

First generate a plan file, then run tfautomv:

```bash
terraform plan -out=tfplan.bin
tfautomv --preplanned
```

#### Single directory with custom plan file

```bash
terraform plan -out=my-plan.bin
tfautomv --preplanned --preplanned-file=my-plan.bin
```

#### Multiple directories

Each directory must have its own plan file:

```bash
# Generate plans in each directory
(cd production && terraform plan -out=tfplan.bin)
(cd staging && terraform plan -out=tfplan.bin)

# Run tfautomv across both directories
tfautomv --preplanned production staging
```

#### JSON vs binary plan files

tfautomv automatically detects the plan file format:
- **Binary plans** (default): tfautomv runs `terraform show -json` to convert them
- **JSON plans** (`.json` extension): tfautomv reads them directly

```bash
# Binary plan (requires terraform show conversion)
terraform plan -out=tfplan.bin
tfautomv --preplanned

# JSON plan (read directly) - useful for CI/CD or when you already have JSON
terraform plan -out=tfplan.bin
terraform show -json tfplan.bin > tfplan.json
tfautomv --preplanned --preplanned-file=tfplan.json
```

**Important**: When using `--preplanned`, all specified directories must have the plan file. If any directory is missing its plan file, tfautomv will exit with an error.

### Disabling colors in output

Add the `--no-color` flag to your `tfautomv` command to disable output
formatting like colors, bold text, etc.

For example:

```bash
tfautomv --no-color
```

Alternatively, you can achieve the same result by setting the `NO_COLOR`
environment variable to any value:

```bash
NO_COLOR=true tfautomv
```

## Thanks

Thanks to [Padok](https://www.padok.fr), where this project was born 💜

## License

The code is licensed under the permissive Apache v2.0 license. [Read this](<https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)>) for a summary.
