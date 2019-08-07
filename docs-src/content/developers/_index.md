---
title: "Developers"
weight: 30
---

{{% children description="true" showhidden="true" depth="999" style="div" %}}

Whether you're developing on Bawt's core functionality or a new plugin, development should always be fun. If you spot something that's not enjoyable then open an issue and let's start fixing it together.

## Prerequisites

- Read the [Code of Conduct](https://github.com/gopherworks/bawt/blob/master/CODE_OF_CONDUCT.md)
- Read our [Contributing Guide](https://github.com/gopherworks/bawt/blob/master/CONTRIBUTING.md)

## Design

Plugins are Go packages in bawt. They are anonymously imported by your main package.

Plugins are represented by structs. At a minimum your plugin struct should include the field `bot *bawt.Bot`.

Plugins take advantage of the `init()` function where you call `bawt.RegisterPlugin(&Help{})`

You'll read more later on about the various plugin interfaces.