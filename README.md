# Automatic Terraform `moved` blocks

When refactoring a Terraform codebase, you often need to write [`moved` blocks](https://www.terraform.io/language/modules/develop/refactoring#moved-block-syntax). This can be tedious. Let
`tfautomv` do it for you.

## How does it work?

When you move a resource in your code, Terraform loses track of the resource's
state. The next time you run Terraform, it will plan to delete the resource it
has memory of and create the "new" resource it found in your refactored code.

`tfautomv` inspects the output of `terraform plan`, detects such
creation/deletion pairs and writes a `moved` block so that Terraform now knows
no deletion or creation is required.

To match resources, `tfautomv` looks for resources of the same type with
identical attributes. If a single match is found, then it writes a `moved`
block.

For more details, see [docs/design.md](./docs/design.md).

## Goals

This project's goals are, in order:

1. `tfautomv` should be certain of the matches it finds
2. `tfautomv` should be quick and painless to use
3. `tfautomv` should be idempotent

## Installation

Compile `tfautomv` from source:

```bash
make build
```

And then put `bin/tfautomv` into your system's PATH.

## Usage

In any directory where you would run `terraform plan`, simply run:

```bash
tfautomv
```

In the background, `tfautomv` will run `terraform plan` and append `moved`
blocks to a `moves.tf` file.

This is what a usual workflow looks like:

```text
$ terraform plan
...
Plan: 19 to add, 0 to change, 19 to destroy.
...

$ tfautomv
Generating moved blocks...
Generated 19 moved blocks.

$ terraform plan
...
Plan: 0 to add, 0 to change, 0 to destroy.
...
```

## Known issues

Some configurations make it so that `tfautomv` cannot be certain about which
resource matches which. For more information, see the failing
`TestDetermineMovedBlocks/dependencies` test in the `internal/tfautomv` package.
