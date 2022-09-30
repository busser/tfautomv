---
weight: 4
title: "Print state mv commands"
description: Tfautomv can provide commands compatible with versions of Terraform older than 1.1.
---

# Print `state mv` commands

Add the `-output=commands` flag to your `tfautomv` command to write
`terraform state mv` commands to standard output instead of writing `moved`
blocks to a file:

```console
$ tfautomv -output=commands
Running "terraform init"...
Running "terraform plan"...
terraform state mv "random_pet.bird" "random_pet.this[\"bird\"]"
terraform state mv "random_pet.cat" "random_pet.this[\"cat\"]"
terraform state mv "random_pet.dog" "random_pet.this[\"dog\"]"
terraform state mv "random_pet.turtle" "random_pet.this[\"turtle\"]"
╷
│ Done: Wrote 4 commands to standard output.
╵
```

The rest of `tfautomv`'s output is written to standard error, not standard
output. This means you can pipe those commands to a file like this:

```console
$ tfautomv -output=commands > moves.sh
Running "terraform init"...
Running "terraform plan"...
╷
│ Done: Wrote 4 commands to standard output.
╵
$ cat moves.sh
terraform state mv "random_pet.bird" "random_pet.this[\"bird\"]"
terraform state mv "random_pet.cat" "random_pet.this[\"cat\"]"
terraform state mv "random_pet.dog" "random_pet.this[\"dog\"]"
terraform state mv "random_pet.turtle" "random_pet.this[\"turtle\"]"
```

Or even pipe them into a shell to run them immediately:

```console
$ tfautomv -output=commands | bash
Running "terraform init"...
Running "terraform plan"...
╷
│ Done: Wrote 4 commands to standard output.
╵
Move "random_pet.bird" to "random_pet.this[\"bird\"]"
Successfully moved 1 object(s).
Move "random_pet.cat" to "random_pet.this[\"cat\"]"
Successfully moved 1 object(s).
Move "random_pet.dog" to "random_pet.this[\"dog\"]"
Successfully moved 1 object(s).
Move "random_pet.turtle" to "random_pet.this[\"turtle\"]"
Successfully moved 1 object(s).
```
