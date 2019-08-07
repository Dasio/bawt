---
title: "Chat Plugins"
weight: 20
---

Chat Plugins use the [PluginInitializer interface](http://localhost:1313/bawt/developers/plugins/types/#plugininitializer).

## Listeners

In bawt one way you can interact with users is over chat. A listener is capable of listening for messages that align with a multitude of criteria.

### Fields

In bawt listeners are structs that are highly configurable via various fields. Let's take a look at the options:

| Field | Type | Description |
| :-- | --- | :-- |
| Name | string | Name of the app. Used during app listing. |
| Description | string | Description of the app. Used during app listing. |
| Slug | string | Slug is a short code used in the help menu |
| Commands | []Command | Commands are the help documentation for commands |
| ListenUntil | time.Time | ListenUntil sets an absolute date at which this Listener expires and stops listening.  ListenUntil and ListenDuration are optional and mutually exclusive. |
| ListenDuration | time.Duration | ListenDuration sets a timeout Duration, after which this Listener stops listening and is garbage collected. A call to `ResetTimeout()` restarts the listening period for another `ListenDuration`. |
| FromUser | *slack.User | FromUser filters out incoming messages that are not with `*User` (publicly or privately)
| FromChannel | *Channel | FromChannel filters messages that are sent to a different room than `Room`. This can be mixed and matched with `FromUser`
| FromAdmin | bool | FromAdmin filters messages that are only meant to be said by an admin |
| FromGroup | []slack.Group | FromGroup filters messages that are not from these groups |
| FromInternalGroup | []string | FromInternalGroup filters out messages not from these groups |
| PrivateOnly | bool | PrivateOnly filters out public messages |
| PublicOnly | bool | PublicOnly filters out private messages.  Mutually exclusive with `PrivateOnly` |
| Contains | string | Contains checks whether the `string` is in the message body (after lower-casing both components) |
| ContainsAny | []string | ContainsAny checks that any one of the specified strings exist as substrings in the message body.  Mutually exclusive with `Contains` |
| Matches | *regexp.Regexp | Matches checks that the given text matches the given Regexp with a `FindStringSubmatch` call. It will set the `Message.Match` attribute. {{% notice note %}} If you spin off a goroutine in the MessageHandlerFunc, make sure to keep a copy of the `Message.Match` object because it will be overwritten by the next Listener the moment your MessageHandlerFunc unblocks {{% /notice %}} |
| ListenForEdits | bool | ListenForEdits will trigger a message when a user edits a message as well as creates a new one |
| MentionsMeOnly | bool | MentionsMe filters out messages that do not mention the Bot's `bot.Config.MentionName` |
| MatchMyMessages | bool | MatchMyMessages equal to false filters out messages that the bot itself sent. |
| MessageHandlerFunc | func(*Listener, *Message) | MessageHandlerFunc is a handling function provided by the user, and called when a relevant message comes in |
| EventHandlerFunc | func(*Listener, interface{}) | EventHandlerFunc is a handling function provided by the user, and called when any event is received. These messages are dispatched to each Listener in turn, after the bot has processed it. If the event is a Message, then the `bawt.Message` will be non-nil. When receiving a `*slack.MessageEvent`, bawt will wrap it in a `*bawt.Message` which embeds the the original event, but adds quite a few functionalities, like reply modes, etc..
| TimeoutFunc | func(*Listener) | TimeoutFunc is called when a conversation expires after `ListenDuration` or `ListenUntil` delays.  It is *not* called if you explicitly call `Close()` on the conversation, or if you did not set `ListenDuration` nor `ListenUntil`. Also, if you override TimeoutFunc, you need to call Close() yourself otherwise, the conversation is not removed from the listeners |
| Bot | *Bot | Bot is a reference to the bot instance.  It will always be populated before being passed to handler functions.

### Methods

Methods extend granular control to developers, allowing them to manually close a listener, control ACK's to Slack's servers, and reset timers.

| Method | Description |
| :-- | :-- |
| Close() | Close terminates the Listener management goroutine, and stops any further listening and message handling |
| ReplyAck() | ReplyAck returns the AckMessage received that corresponds to the Reply on which you called `Listen()` |
| ResetDuration() | ResetDuration re-initializes the timeout set by `Listener.ListenDuration`, and continues listening for another such duration. |

## Message Handling

When you receive a message after it matches the criteria given by a `bawt.Listener` you will receive it as a struct called `bawt.Message`.

### Fields

| Field | Type | Description |
| :-- | :-- | :-- |
| `*slack.Msg` | | Allows you to accesss the underlying Slack message primitive |
| `SubMessage` | `*slack.Msg` | Threaded replies |
| `MentionsMe` | `bool` | Does the user @mention the bot |
| `IsEdit` | `bool` | Is this message the result of an edit |
| `FromMe` | `bool` | Was the message from the bot |
| `FromUser` | `*slack.User` | The user who sent the message |
| `FromChannel` | `*Channel` | The channel the message was received in, including direct messages |
| `Match` | `[]string` | Match contains the result of Listener.Matches.FindStringSubmatch(msg.Text), when `Matches` is set on the `Listener`. |

### Methods

| Function Signature | Description |
| :-- | :-- |
| `AddReaction(emoticon string) *Message` | AddReaction adds a reaction to a message |
| `Contains(s string) bool` | Contains searches for a single string in a noncase-sensitive fashion |
| `ContainsAll(strs []string) bool` | ContainsAll searches for all strings in a noncase-sensitive fashion |
| `ContainsAny(strs []string) bool` | ContainsAny searches for at least one noncase-sensitive matching string |
| `ContainsAnyCased(strs []string) bool` | ContainsAnyCased searches for at least one case-sensitive word |
| `HasPrefix(prefix string) bool` | HasPrefix returns true if a message starts with a given string |
| `IsPrivate() bool` | IsPrivate determines if a message is private or not |
| `ListenReaction(reactListen *ReactionListener)` | ListenReaction listens for a reaction on a message |
| `RemoveReaction(emoticon string) *Message` | RemoveReaction removes a reaction from a message |
| `Reply(text string, v ...interface{}) *Reply` | Reply sends a message back to the source it came from, without a mention |
| `ReplyMention(text string, v ...interface{}) *Reply` | ReplyMention replies with a @mention named prefixed, when replying in public. When replying in private, nothing is added. |
| `ReplyPrivately(text string, v ...interface{}) *Reply` | ReplyPrivately replies to the user in an IM |
| `ReplyWithFile(p FileUploadParameters) *ReplyWithFile` | ReplyWithFile replies with a snippet or an attached file |
| `String() string` | String returns a message with field:value as a string |