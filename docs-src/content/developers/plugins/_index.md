---
title: "Plugins"
weight: 10
---

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

## Interfaces

Bawt's plugins must only satisfy an `interface`, which in Go means that we must simply adhere to a certain minimal method structure.

We do some polymorphism here, but it's not really complicated, so stay with us!

### The base interface

The base interface is `Plugin` which is satisfied by a struct with no methods.

### The chat interface

The `PluginInitializer` interface is satisifed by a method of `InitPlugin(*Bot)`. In the context of our `Help` plugin, let's look how we'd use that!

```go
// InitPlugin initializes the plugin
func (h *Help) InitPlugin(bot *bawt.Bot) {
	h.bot = bot

	h.listenHelp()
}

func (h *Help) listenHelp() {
	h.bot.Listen(&bawt.Listener{
		Matches:            regexp.MustCompile(`^!help.*`),
		MessageHandlerFunc: h.handleHelp,
	})
}

func (h *Help) handleHelp(listen *bawt.Listener, msg *bawt.Message) {
	// Do some things
	msg.Reply("Got your message!")
}
```

Messages are sent to registered listeners as they are received and if the listener regex matches then we execute the `MessageHandlerFunc` which takes two arguments supplied by Bawt's core: `Listener` and `Message`.

While `Listener` has some useful exported methods, the `Message` type comes with methods for replying in various ways and will be what a majority of developers will be interacting with.

### The webserver interface

Bawt can extend itself to listen to web requests.

The web server plugin interface is fully functional, however, not fully documented. It is currently under development and will receive a lot better documentation in the future :)