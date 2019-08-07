---
title: "Plugins"
weight: 10
---

{{% children description="true" showhidden="true" depth="999" style="div" %}}

All of Bawt's plugins start with a single `struct`. This struct stores configuration data that will be unmarshaled to it during the init process.

At a minimum a field called `bot` which is a pointer of type `bawt.Bot`. This field is required in order to create your first listener, so don't forget it!

```go
type Help struct {
	bot *bawt.Bot
}
```

We'll then use Go's built in `init()` function to register the plugin at runtime:


```go
func init() {
	bawt.RegisterPlugin(&Help{})
}
```

Now our plugin is registered, but it still needs to be initialized.

On the next page we'll be learning about the different plugin types and the methods they support!