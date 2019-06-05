---
title: "Getting Started"
weight: 10
---

Bawt is a chatops _framework_ rather than a bot in itself. This means we don't provide things like the scaffolding for a command line but it also means we remain fairly extensible. This structure allows teams to run their own Bawt without every Bawt being the same.

Using Bawt is so simple that two files is all it takes! Our mission is to make Bawt easy and fun to use for everyone, so follow along and let's connect our first bot!

## Your First Bot

_This tutorial will assume you are working out of your home directory (`~` or `$HOME`)_

1\. Create a new directory for our bot:

  - `mkdir new-bot && cd new-bot`

2\. New-Bot will need a binary, so let's create a minimal `main.go`. The contents should resemble this:

```go
package main

import (
	"flag"

	"github.com/gopherworks/bawt"
	_ "github.com/gopherworks/bawt/help"
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

```

- Turn your eyes to the imports for just a second. Bawt's code takes advantage of imports with _blank identifiers_ and the package `init()` functions to load plugins. You do not need to download plugins, as long as you obtain them using modules they'll be versioned for you as well!

3\. Almost there! Let's tell New-Bot how to connect to Slack. Reading through our code earlier we're looking for a config file. Bawt will always look for `config.(json|yaml|toml)` in the current directory and your home directory by default, but `bawt.New()` takes an optional pointer to override this file.

```json
{
  "Slack": {
    "api_token": "xoxb-mamamamama-papapapapapapapa",
    "nickname": "username",
    "general_channel": "#general",
    "team_domain": "your-team-domain-name",
    "join_channels": [
      "#some", "#other", "private_group"
    ],
    "web_base_url": "http://host.example.com",
    "db_path": "./bawt.bolt.db"
  },

  "Server":{
    "pid_file": "/var/run/bawt.pid-or-empty-string"
  },
}
```

- Each plugin may require different configuration entries but `Server` and `Slack` are the two required entries.

4\. To make things easy, let's build a Makefile!

```Makefile
clean:
	@printf "# Removing vendor dir\n"
	@rm -rf vendor
	@printf "# Removing build dir\n"
	@rm -rf build

build: clean vendor
	@env GO111MODULE=on go build -o build/new-bot .
	@chmod a+x build/new-bot

vendor:
	@go mod tidy
	@go mod vendor

run:
    ./bawt/new-bot --config config.json
```

That's it! We're ready to fire New-Bot up!

`make run`