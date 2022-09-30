---
title: Automatic Terraform moved blocks
type: docs
description: Easier refactoring, less technical debt.
---

# Automatic Terraform `moved` blocks

{{< columns >}}

## Easier refactoring

Terraform's inherent statefulness makes it painful to refactor an existing
codebase. Tfautomv writes `moved` blocks for you so your refactoring is quicker
and less error-prone.

<--->

## Less technical debt

Refactoring prevents your code drifting away from the mental model you have
of your infrastructure. We call this drift _technical debt_. Having less of it
makes you more productive.

{{< /columns >}}

We explain why we built tfautomv in more detail [in this blog article](https://www.padok.fr/en/blog/terraform-refactoring-tfautomv).

## Getting started

Start with our [quick-start guide]({{< relref "getting-started/_index.md" >}}) to set up tfautomv on your
workstation.

Once you are set up, learn more about [tfautomv's features]({{< relref "usage/_index.md" >}}).

## Going further

To understand how tfautomv works under the hood, [read our design]({{< relref "design/_index.md" >}}).
