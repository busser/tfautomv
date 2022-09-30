---
weight: 3
title: "Design"
---

# Design

## Goals

This project's goals are, in order:

1. `tfautomv` should be certain of the matches it finds
2. `tfautomv` should be quick and painless to use
3. `tfautomv` should be idempotent

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

## Assumptions

1. The changes to the codebase do not require any changes to the managed
   resources. Once the correct `moved` blocks have been generated, running
   `terraform plan` should yield no planned changes.
2. Resources are moved around in the same Terraform state. `tfautomv` works
   based on the output of the `terraform plan` command, so working across
   multiple states is not in the project's scope.

## Examples

Here are some examples of what we want `tfautomv` to be able to do.

### A single resource of a given type

If we have made the following changes to our codebase:

```diff
- resource "random_id" "foo" {
+ resource "random_id" "bar" {
    byte_length = 6
  }
```

then we want `tfautomv` to generate this code:

```terraform
moved {
  from = random_id.foo
  to   = random_id.bar
}
```

### Multiple resources of the same type

If we have made changes to multiple resources of the same type:

```diff
- resource "random_id" "first" {
+ resource "random_id" "alpha" {
    byte_length = 6
  }

- resource "random_id" "second" {
+ resource "random_id" "beta" {
    byte_length = 8
  }
```

then we want `tfautomv` to generate this code:

```terraform
moved {
  from = random_id.first
  to   = random_id.alpha
}

moved {
  from = random_id.second
  to   = random_id.beta
}
```

We want `tfautomv` to use the differences between the resources known attributes
to match resources 1 to 1.

Using differences between resource attributes has limitations. For example, if
we make the following changes to the codebase:

```diff
- resource "random_id" "first" {
+ resource "random_id" "alpha" {
    byte_length = 6
  }

- resource "random_id" "second" {
+ resource "random_id" "beta" {
    byte_length = 6
  }
```

We want `tfautomv` to generate the same `moved` blocks as above. However, the
two `random_id` resources have identical values for all their known attributes.
This means that `tfautomv` could just as well generate the following code:

```terraform
moved {
  from = random_id.first
  to   = random_id.beta
}

moved {
  from = random_id.second
  to   = random_id.alpha
}
```

This mapping could be wrong and might force unnecessary changes to resources
that depend on the `random_id` resources.

### Looking at dependencies

Let's extend the previous example:

```diff
- resource "random_id" "first" {
+ resource "random_id" "alpha" {
    byte_length = 6
  }

- resource "random_id" "second" {
+ resource "random_id" "beta" {
    byte_length = 6
  }

  resource "google_sql_database_instance" "alpha" {
-   name = "alpha-${random_id.first.hex}"
+   name = "alpha-${random_id.alpha.hex}"

    database_version = "MYSQL_5_7"

    // ...
  }

  resource "google_sql_database_instance" "beta" {
-   name = "beta-${random_id.second.hex}"
+   name = "beta-${random_id.beta.hex}"

    database_version = "POSTGRES_11"

    // ...
  }
```

By looking only at the `random_id` resources, we cannot make a deterministic
mapping. However, the `google_sql_database_instance` resources can be
distinguished based solely on theur known attributes (ie. `database_version`).

A Terraform plan includes information about resource dependencies. By analysing
those dependencies, `tfautomv` should be able to link each `random_id` resource
to a single `google_sql_database_instance` resource. A set of `moved` blocks can
then be generated for each group of linked resources, yielding the correct code:

```terraform
moved {
  from = random_id.first
  to   = random_id.alpha
}

moved {
  from = random_id.second
  to   = random_id.beta
}
```
