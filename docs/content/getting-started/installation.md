---
weight: 2
title: "Installation"
description: How to set up tfautomv on your workstation.
---

# Installation

Follow any of the instructions below.

## Mac OS

### Pre-compiled binary

```bash
brew install busser/tap/tfautomv
```

_This formula is maintained by the core tfautomv team. It is updated
automatically after every release._

## General Linux

### Pre-compiled binary

```bash
curl -L https://raw.githubusercontent.com/busser/tfautomv/main/install.sh | bash
```

_This script installs the latest release by default._

## Arch Linux

You can install tfautomv with the [Arch User Repository](https://wiki.archlinux.org/title/Arch_User_Repository).

### Pre-compiled binary

```bash
yay tfautomv-bin
```

_The [`tfautomv-bin` AUR package](https://aur.archlinux.org/packages/tfautomv-bin)
is maintained by the core tfautomv team. It is updated automatically after every
release._

### From source

```bash
yay tfautomv
```

_The [`tfautomv` AUR package](https://aur.archlinux.org/packages/tfautomv) is
maintained by a member of the community._

## Manual

{{< hint info >}}

The tfautomv team is [open to requests](https://github.com/busser/tfautomv/issues)
of other installation methods.

{{< /hint >}}

### Pre-compiled binary

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

## Next steps

Confirm that tfautomv is properly installed:

```bash
tfautomv -version
```

Then Follow the [guided tutorial]({{< relref "getting-started/tutorial.md" >}})
to become familiar with tfautomv.
