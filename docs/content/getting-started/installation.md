---
weight: 2
title: "Installation"
description: How to set up tfautomv on your workstation.
---

# Installation

Follow the instructions in any of the tabs below:

{{< tabs "installation" >}}
{{< tab "Homebrew" >}}

```bash
brew install padok-team/tap/tfautomv
```

{{< /tab >}}
{{< tab "Pre-compiled binary" >}}

On the Github reposotory's [Releases page](https://github.com/padok-team/tfautomv/releases),
download the binary that matches your workstation's OS and CPU architecture.

Put the binary in a directory present in your system's `PATH` environment
variable.

{{< /tab >}}
{{< tab "From source" >}}

You must have Go 1.18+ installed to compile tfautomv.

Clone the repository and build the binary:

```bash
git clone https://github.com/padok-team/tfautomv
cd tfautomv
make build
```

Then, move `bin/tfautomv` to a directory resent in your system's `PATH`
environment variable.

{{< /tab >}}
{{< tab "Distro Packages" >}}
Available in AUR for Arch Linux: https://aur.archlinux.org/packages/tfautomv

```
yay tfautomv
```

{{< /tab >}}
{{< /tabs >}}

Confirm that tfautomv is properly installed:

```bash
tfautomv -version
```

## Next steps

Follow the [guided tutorial]({{< relref "getting-started/tutorial.md" >}}) to become familiar with tfautomv.
