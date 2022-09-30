---
weight: 2
title: "Ignore certain differences"
description: Tfautomv can ignore differences between certain resources' attributes, based on rules you provide.
---

# Ignore certain differences

Add the `-ignore` flag to your `tfautomv` command to ignore differences between
certain resources' attributes. Differences are ignored based on the rules you
provide. Use the `-ignore` flag multiple times to specify multiple rules.

```plaintext
<EFFECT>:<RESOURCE TYPE>:<ATTRIBUTE NAME>
```

For example:

```bash
tfautomv -ignore=everything:random_pet:length
```

For nested attributes, separate parent attributes from child attributes with a
`.` (the representation used in tfautomv's detailed analysis):

```bash
<EFFECT>:<RESOURCE TYPE>:parent_obj.child_field
<EFFECT>:<RESOURCE TYPE>:parent_list.0
```

## Available effects

The following effects are available:

- `everything`: ignores all differences between attribute values
- `whitespace`: ignores whitespace when comparing attribute values
