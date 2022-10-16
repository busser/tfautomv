---
weight: 3
title: "Tutorial"
description: Using tfautomv for the first time.
---

# Tutorial

## 0. Requirements

Before starting this tutorial, you must:

- Install Terraform
- [Install tfautomv]({{< relref "getting-started/installation.md" >}})

Both these commands should work:

```bash
terraform version
tfautomv -version
```

## 1. Initial Terraform code

Start by creating an empty directory. This will be your workspace for the
duration of this tutorial. Let's write some Terraform code into it.

In order to make this tutorial as accessible as possible, we won't be
provisioning any actual resources. Instead, we will use Terraform built-in
`random` provider because it does not require credentials of any kind.

{{< hint info >}}

`tfautomv` is fully **provider-agnostic.**
It works with all Terraform providers.

{{< /hint >}}

Create a single file called `main.tf` with the following content:

```terraform
resource "random_pet" "bird" {
  prefix = "bird"
  length = 3
}

resource "random_pet" "cat" {
  prefix = "cat"
  length = 3
}

resource "random_pet" "dog" {
  prefix = "dog"
  length = 3
}

resource "random_pet" "turtle" {
  prefix = "turtle"
  length = 3
}
```

Run Terraform to create the four resources:

```bash
terraform init
terraform apply
```

Terraform plans to create the resources, as expected:

```console
Plan: 4 to add, 0 to change, 0 to destroy.
```

Once approved, Terraform puts its plan into action:

```console
Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

If you run Terraform again, it says that no changes are required:

```console
No changes. Your infrastructure matches the configuration.
```

This makes sense, because the resources already exist.

Now let's refactor that codebase.

## 2. Refactoring

There are many ways to structure a Terraform codebase. Of all the different ways
you could structure your code, which one should you choose? Our advice: write
your code so that it matches the mental model you have of your infrastructure.

{{< hint info >}}

Teams usually have two models of their infrastructure. The first is their code,
which matches the infrastructure as it exists in the cloud. The second is the
mental model of their infrastructure. When these two models drift — when a
team's mental model differs from their actual code — we call this
**technical debt**.

{{< /hint >}}

For this tutorial, let's say that in our mental model these `random_pet`
resources all have a unique prefix. This fact appears nowhere in our code. Let's
edit our code to reflect that these prefixes are unique. Update `main.tf` with
this new content:

```terraform
locals {
  pets = toset([
    "bird",
    "cat",
    "dog",
    "turtle"
  ])
}

resource "random_pet" "this" {
  for_each = local.pets

  prefix = each.key
  length = 3
}
```

Now your code reflects the fact that each `random_pet` has a unique prefix,
because of the `toset` function which would remove any duplicates.

See what Terraform plans to do following these code changes:

```bash
terraform plan
```

You should see this scary message:

```console
Plan: 4 to add, 0 to change, 4 to destroy.
```

But that is not what you want. In a real-world scenario, Terraform may plan to
destroy production resources which is most likely unacceptable. Why does
Terraform not understand these are the same resources?

Terraform keeps track of the existing resources in its state. Currently, your
resources are in Terraform's state with these ID's:

```console
$ terraform state list
random_pet.bird
random_pet.cat
random_pet.dog
random_pet.turtle
```

But your code no longer uses those ID's. It uses these instead:

```plaintext
random_pet.this["bird"]
random_pet.this["cat"]
random_pet.this["dog"]
random_pet.this["turtle"]
```

Since the resources Terraform knows about are no longer in your code, Terraform
wants to destroy them. Inversely, since Terraform does not know about the
resources in your code, it assumes it must create them.

So what should you do?

## 3. Use tfautomv

Using [`moved` blocks](https://www.terraform.io/language/modules/develop/refactoring#moved-block-syntax), `tfautomv` can tell Terraform how to update its
state to match your new code. All you need to do is run it:

```bash
tfautomv
```

You should see this output:

```console
Running "terraform init"...
Running "terraform plan"...
╷
│ Done: Added 4 moved blocks to "moves.tf".
╵
```

Now run Terraform again to see what it has planned:

```bash
terraform apply
```

You should see this message:

```console
Plan: 0 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
```

All is good, Terraform does not plan to change any resources! The prompt is
asking you if you want to enact the moves `tfautomv` wrote for you. The moves
are correct, so say yes.

Once you see this message, you have finished the tutorial:

```console
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

## Next steps

Learn more about [tfautomv's features]({{< relref "usage/_index.md" >}}) or about
[how it works]({{< relref "design/_index.md" >}}) under the hood.
