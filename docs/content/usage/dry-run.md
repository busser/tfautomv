---
weight: 1
title: "Perform a dry run"
description: Tfautomv can analyse Terraform's plan without making any changes to your code.
---

# Perform a dry run

Add the `-dry-run` flag to your `tfautomv` command to print moves in a
human-readable format instead of writing them to a file:

```console
$ tfautomv -dry-run
Running "terraform init"...
Running "terraform plan"...
╷
│ Moves
│ ╷
│ │ From: random_pet.bird
│ │ To:   random_pet.this["bird"]
│ ╵
│ ╷
│ │ From: random_pet.cat
│ │ To:   random_pet.this["cat"]
│ ╵
│ ╷
│ │ From: random_pet.dog
│ │ To:   random_pet.this["dog"]
│ ╵
│ ╷
│ │ From: random_pet.turtle
│ │ To:   random_pet.this["turtle"]
│ ╵
╵
```
