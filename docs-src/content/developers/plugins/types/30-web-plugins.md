---
title: "Web Plugins"
weight: 30
---

Web plugins primarily rely on `gorilla/mux` to serve HTTP or HTTPS pages. Most of the setup of a web plugin is done inside the `InitWebPlugin` function. You're provided a reference to the bot, a private router, and a public router. Let's look at the example below:

```go
func (wp *WebPlugin) InitWebPlugin(bot *bawt.Bot, privRouter *mux.Router, pubRouter *mux.Router) {
    // Storing the bawt reference
    wp.bot = bot

    // Load some config; see: https://gopherworks.github.io/bawt/developers/plugins/20-useful-functions/ 
    var conf struct {
        wp wpConfig
    }
    bot.LoadConfig(&conf)
    wp.config = conf.wp

    // The public router uses /public as a prefix
    pubRouter.HandleFunc("/public/ping", wp.handlePing)

    // The private router runs on localhost
    privRouter.HandleFunc("/healthz", wp.handleHealthcheck)
}
```

As you can see we load some configuration using bawt's built in functionality for loading it's core configuration. Next we setup a route on the public router for `/public/ping` and send it to a method called `handlePing`. This router is served on the IP address specified in the configuration.

Private routers are just the same except they serve on `localhost`.