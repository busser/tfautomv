---
weight: 2
title: "Ignore certain differences"
description: Tfautomv can ignore differences between certain resources' attributes, based on rules you provide.
---

# Ignore certain differences

Add the `-ignore` flag to your `tfautomv` command to ignore differences between
certain resources' attributes. Differences are ignored based on the rules you
provide. Use the `-ignore` flag multiple times to specify multiple rules.

Rules have this structure:

```plaintext
<EFFECT>:<RESOURCE TYPE>:<ATTRIBUTE NAME>[OPTIONAL PARAMETERS]
```

For nested attributes, separate parent attributes from child attributes with a
`.` (the representation used in tfautomv's [detailed analysis]({{< relref "usage/show-analysis.md" >}})):

```bash
<EFFECT>:<RESOURCE TYPE>:parent_obj.child_field
<EFFECT>:<RESOURCE TYPE>:parent_list.0
```

## Ignore an attribute entirely

Use the `everything` effect to ignore any difference between two values of an
attribute:

```bash
tfautomv -ignore="everything:<RESOURCE TYPE>:<ATTRIBUTE NAME>"
```

For example:

```bash
tfautomv -ignore="everything:random_pet:length"
```

## Ignore whitespace

Use the `whitespace` effect to ignore differences in whitespace between two
values of an attribute:

```bash
tfautomv -ignore="whitespace:<RESOURCE TYPE>:<ATTRIBUTE NAME>"
```

For example:

```bash
tfautomv -ignore="whitespace:azurerm_api_management_policy:xml_content"
```

## Ignore a prefix

Use the `prefix` effect to ignore a specific prefix between in one of two values
of an attribute:

```bash
tfautomv -ignore="prefix:<RESOURCE TYPE>:<ATTRIBUTE NAME>:<PREFIX>"
```

For example:

```bash
tfautomv -ignore="prefix:google_storage_bucket_iam_member:bucket:b/"
```
