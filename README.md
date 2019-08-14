![Build Status](https://travis-ci.com/gopherworks/bawt.svg?branch=master) 
![Release](https://img.shields.io/github/release/gopherworks/bawt.svg)
[![GoDoc](https://godoc.org/github.com/gopherworks/bawt?status.svg)](https://godoc.org/github.com/gopherworks/bawt) 
[![Go Report Card](https://goreportcard.com/badge/github.com/gopherworks/bawt)](https://goreportcard.com/report/github.com/gopherworks/bawt)
![License](https://img.shields.io/github/license/gopherworks/bawt.svg)
![Open Issues](https://img.shields.io/github/issues-raw/gopherworks/bawt.svg)
![Open PRs](https://img.shields.io/github/issues-pr-raw/gopherworks/bawt.svg)

Bawt
===

Bawt is a chatops _framework_ rather than a bot in itself, hence, Bawt is distributed as a package and all it takes to start is one file and one function.

Our goal is that bawt is always **easy to start**, **easy to run**, **easy to enhance**

**Easy to start.** The single biggest turn off I had when trying out a new project was the amount of time it took me to get started with something meaningful. The plugin footprint is intentionally light yet extensive and why Bawt's core can be started with just two files. We aim to strike a rare balance of extensibility and simplicity.

**Easy to run.** Bawt is updated via modules and follows [Semantic Versioning (SemVer)](https://semver.org/) so you'll always know what sort of changes await you. Bawt's core code is abstracted into the Messaging API so even when Slack breaks their API's (and they will) you will never notice as a plugin developer.

**Easy to enhance.** We want Bawt's code to make sense not just to core developers but to plugin devs as well. That's why Bawt's core is both verbose and descriptive. There's no buried functionality, what you see is what you get.

## Local build and install

### Using Bawt

```shell
go get github.com/gopherworks/bawt
```

You import Bawt like any other package. Learn about [getting started with Bawt](https://gopherworks.github.io/bawt/getting-started/)!

## Writing your own plugin

Example code to handle deployments:

```go
// listenDeploy was hooked into a plugin elsewhere..
func listenDeploy() {
	keywords := []string{"project1", "project2", "project3"}
	bot.Listen(&bawt.Listener{
		Matches:        regexp.MustCompile("(can you|could you|please|plz|c'mon|icanhaz) deploy (" + strings.Join(keywords, "|") + ") (with|using)( revision| commit)? `?([a-z0-9]{4,42})`?"),
		MentionsMeOnly: true,
		MessageHandlerFunc: func(listen *bawt.Listener, msg *bawt.Message) {

			projectName := msg.Match[2]
			revision := msg.Match[5]

			go func() {
				go msg.AddReaction("work_hard")
				defer msg.RemoveReaction("work_hard")

				// Do the deployment with projectName and revision...

			}()
		},
	})
}
```

Whether you're interested in developing on Bawt's core or your own plugin our [developer docs](https://gopherworks.github.io/bawt/developers/) can help out.