# tfautomv multi-state test fixture

This fixture contains three independent Terraform root modules with local
state. It uses Terraform's built-in `terraform_data` resource, so it does not
need cloud credentials or an external provider.

## Create the baseline states

From this directory, initialize and apply each stack independently:

```bash
for stack in network shared application; do
  terraform -chdir="$stack" init
  terraform -chdir="$stack" apply
done
```

Each stack has multiple `.tf` files and several resources. The apply commands
create only local Terraform state; `terraform_data` does not create external
infrastructure.

## Move resources between stacks

For a cross-state test, move one complete resource declaration by copying its
file and removing the original. For example:

```bash
cp network/gateway.tf application/gateway.tf
rm network/gateway.tf
```

Create fresh plans after the configuration change:

```bash
for stack in network shared application; do
  terraform -chdir="$stack" plan -out=tfplan.bin
done
```

From the repository root, inspect the mapping without writing anything:

```bash
tfautomv --preplanned --dry-run \
  examples/multi-state-test/network \
  examples/multi-state-test/shared \
  examples/multi-state-test/application
```

Require every planned create and delete to be paired:

```bash
tfautomv --preplanned --dry-run --require-net-zero \
  examples/multi-state-test/network \
  examples/multi-state-test/shared \
  examples/multi-state-test/application
```

## Generate migration files

Generate `moved` blocks for same-state moves and `removed`/`import` blocks in
the appropriate source and destination directories:

```bash
tfautomv --preplanned --output=blocks \
  examples/multi-state-test/network \
  examples/multi-state-test/shared \
  examples/multi-state-test/application
```

Review the generated `moves.tf` or `state-migration.tf` files, then regenerate
plans and apply them in the affected stacks. Remove the temporary migration
blocks after Terraform has recorded the new state addresses.

For a command-based cross-state test, emit destination-first imports followed
by source-state removals:

```bash
tfautomv --preplanned --output=import-commands \
  examples/multi-state-test/network \
  examples/multi-state-test/shared \
  examples/multi-state-test/application
```

The existing `--output=commands` mode remains available for the pulled-local-
state `terraform state mv` workflow.
