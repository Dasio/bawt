package main

import (
	"flag"

	"github.com/gopherworks/bawt"
	_ "github.com/gopherworks/bawt/bugger"
	_ "github.com/gopherworks/bawt/faceoff"
	_ "github.com/gopherworks/bawt/funny"
	_ "github.com/gopherworks/bawt/healthy"
	_ "github.com/gopherworks/bawt/help"
	_ "github.com/gopherworks/bawt/hooker"
	_ "github.com/gopherworks/bawt/mooder"
	_ "github.com/gopherworks/bawt/plotberry"
	_ "github.com/gopherworks/bawt/recognition"
	_ "github.com/gopherworks/bawt/standup"
	_ "github.com/gopherworks/bawt/todo"
	_ "github.com/gopherworks/bawt/web"
	_ "github.com/gopherworks/bawt/webauth"
	_ "github.com/gopherworks/bawt/webutils"
	_ "github.com/gopherworks/bawt/wicked"
)

// Specify an alternative config file. bawt searches the working
// directory and your home folder by default for a file called
// `config.json`, `config.yaml`, or `config.toml` if no config
// file is specified
var configFile = flag.String("config", "", "config file")

func main() {
	flag.Parse()

	bot := bawt.New(*configFile)

	bot.Run()
}
