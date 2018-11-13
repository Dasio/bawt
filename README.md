![Build Status](https://travis-ci.com/gopherworks/bawt.svg?branch=master)

# bawt - A Slack bot in Go

bawt is a Slack bot to do ChatOps and other cool things. bawt was originally developed [here](https://github.com/abourget/slick) and has since been adopted  and sponsored by Capstone :)


## Features

Supported features:

* Plugin interface for chat messages
* Plugin-based HTTP handlers
* Simple API to reply to users
* Keeps an internal state of channels, users and their state.
* Listen for Reactions; take actions based on them (like buttons).
* Simple API to message users privately
* Simple API to update a previously sent message
* Simple API to delete bot messages after a given time duration.
* Easy plugin interface, listeners with criteria such as:
  * Messages directed to the bot only
  * Private or public messages
  * Listens for a duration or until a given `time.Time`
  * Selectively on a channel, or from a user
  * Expire listeners and unregister them dynamically
  * Supports listening for edits or not
  * Regexp match messages, or Contains checks
* Built-in KV store for data persistence (backed by BoltDB and JSON serialization)
* The bot has a mood (_happy_ and _hyper_) which changes randomly.. you can base some decisions on it, to spice up conversations.
* Supports listening for any Slack events (ChannelCreated, ChannelJoined, EmojiChanged, FileShared, GroupArchived, etc..)
* A PubSub system to facilitate inter-plugins (or chat-to-web) communications.


## Stock plugins

1. Recognition: a plugin to recognize your peers (!recognize @user1 for doing an awesome job)

2. Faceoff: a game to learn the names and faces of your colleagues. The code for this one is interesting to learn to build interactive features with `bawt`.

3. Vote: a simple voting plugin to decide where to lunch

4. Funny: a bunch of jokes and memes in reply to some strings in channels.. (inspired by Hubot's jokes)

5. Healthy: a very simple plugin that pokes URLs and reports on their health

6. Deployer: an example plugin to do deployments wth ansible (you'll probably want to roll out your own though).

7. Todo: todo list manager, one per channel


## Local build and install

Try it with:

```shell
go get github.com/gopherworks/bawt
cd $GOPATH/src/github.com/gopherworks/bawt/example-bot
go install -v && $GOPATH/bin/example-bot
```

There's a Dockerfile and example configuration in the `example-bot` directory.

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




Take inspiration by looking at the different plugins, like `Funny`,
`Healthy`, `Storm`, `Deployer`, etc..  Don't forget to update your
bot's plugins list, like in `example-bot/main.go`
