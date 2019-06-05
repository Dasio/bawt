---
title: home
---

## Bawt is a Slack bot framework written in Go. 

---

Our goal is that bawt is always **easy to start**, **easy to run**, **easy to enhance**

**Easy to start.** The single biggest turn off I had when trying out a new project was the amount of time it took me to get started with something meaningful. The plugin footprint is intentionally light yet extensive and why Bawt's core can be started with just two files. We aim to strike a rare balance of extensibility and simplicity.

**Easy to run.** Bawt is updated via modules and follows [Semantic Versioning (SemVer)](https://semver.org/) so you'll always know what sort of changes await you. Bawt's core code is abstracted into the Messaging API so even when Slack breaks their API's (and they will) you will never notice as a plugin developer.

**Easy to enhance.** We want Bawt's code to make sense not just to core developers but to plugin devs as well. That's why Bawt's core is both verbose and descriptive. There's no buried functionality, what you see is what you get.

---

## Features
This is by far not a comprehensive list

* Easy to start, Easy to maintain, Easy to enhance
* Plugin API
  * Single function registration
  * Built in help docs interface
  * Chat, HTTP, or HTTPAuth type plugins
* Messaging API
  * Channel and DM's
  * Public and private messages
  * Ephemeral messages (disappear after duration)
  * Update previously sent message
  * File and Snippet uploads
* Listener API
  * Dynamic registration and deregistration
  * Listen for Messages, Edits, Reactions...
  * Regex and contains based checks
* BoltDB for data persistence
* Central Slack event loop
* Tracks state internally (Channels, Users, and their state)
* Built in web server that support authentication
  * Communication between web plugins and chat plugins
* Uses modules for versioning