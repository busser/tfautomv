# Automatic Terraform `moved` blocks <!-- omit in toc -->

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/padok-team/tfautomv)](https://goreportcard.com/report/github.com/padok-team/tfautomv)
![tests-passing](https://github.com/padok-team/tfautomv/actions/workflows/ci.yml/badge.svg)

When refactoring a Terraform codebase, you often need to write [`moved` blocks](https://www.terraform.io/language/modules/develop/refactoring#moved-block-syntax). This can be tedious. Let
`tfautomv` do it for you.

- [Why does this exist?](#why-does-this-exist)
- [Demo](#demo)
- [Installation](#installation)
- [Usage](#usage)
- [How it works](#how-it-works)
- [License](#license)

## Why does this exist?

When you move a resource in your code, Terraform loses track of the resource's
state. The next time you run Terraform, it will plan to delete the resource it
has memory of and create the "new" resource it found in your refactored code.

`tfautomv` inspects the output of `terraform plan`, detects such
creation/deletion pairs and writes a `moved` block so that Terraform now knows
no deletion or creation is required.

## Demo

This demo illustrates tfautomv's core features:

- Generating `moved` blocks automatically for refactored code
- Optionally doing a dry run
- Optionally showing a detailed analysis

![demo](./docs/content/getting-started/demo.gif)

## Installation

See [Getting started / Installation](https://padok-team.github.io/tfautomv/getting-started/installation/)
for instructions.

## Usage

See [Getting started / Tutorial](https://padok-team.github.io/tfautomv/getting-started/tutorial/)
for a hands-on guided introduction de tfautomv.

See [Usage](https://padok-team.github.io/tfautomv/usage/) for a list of
tfautomv's features.

The following versions of Terraform are supported:

- `1.1.x` and above by default
- `0.13.x` and above when using the `-output=commands` flag

## How it works

See [Design](https://padok-team.github.io/tfautomv/design/) for details on how tfautomv works under the hood.

## License

The code is licensed under the permissive Apache v2.0 license. [Read this](<https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)>) for a summary.
