package bawt

import (
	"time"

	"github.com/nlopes/slack"
)

// ReactionListener listens for reactions and changes in reactions
type ReactionListener struct {
	ListenUntil    time.Time
	ListenDuration time.Duration
	FromUser       *slack.User
	Emoji          string
	Type           reaction

	HandlerFunc func(listen *ReactionListener, event *ReactionEvent)
	TimeoutFunc func(*ReactionListener)

	listener *Listener
}

func (rl *ReactionListener) newListener() *Listener {
	newListen := &Listener{}
	if !rl.ListenUntil.IsZero() {
		newListen.ListenUntil = rl.ListenUntil
	}
	if rl.ListenDuration != time.Duration(0) {
		newListen.ListenDuration = rl.ListenDuration
	}
	if rl.TimeoutFunc != nil {
		newListen.TimeoutFunc = func(listen *Listener) {
			rl.TimeoutFunc(rl)
		}
	}
	rl.listener = newListen

	return newListen
}

func (rl *ReactionListener) filterReaction(re *ReactionEvent) bool {
	if rl.Emoji != "" && re.Emoji != rl.Emoji {
		return false
	}
	if rl.FromUser != nil && re.User != rl.FromUser.ID {
		return false
	}
	if int(rl.Type) != 0 && re.Type != rl.Type {
		return false
	}
	return true
}

// Close closes the connection
func (rl *ReactionListener) Close() {
	rl.listener.Close()
}

// ResetNewDuration resets the duration timer and creates a new duration
func (rl *ReactionListener) ResetNewDuration(d time.Duration) {
	rl.listener.ListenDuration = d
	rl.listener.ResetDuration()
}

// ResetDuration resets the duration timer
func (rl *ReactionListener) ResetDuration() {
	rl.listener.ResetDuration()
}

// ReactionEvent is a reaction event from Slack
type ReactionEvent struct {
	// Type can be `ReactionAdded` or `ReactionRemoved`
	Type      reaction
	User      string
	Emoji     string
	Timestamp string
	Item      struct {
		Type        string `json:"type"`
		Channel     string `json:"channel,omitempty"`
		File        string `json:"file,omitempty"`
		FileComment string `json:"file_comment,omitempty"`
		Timestamp   string `json:"ts,omitempty"`
	}

	// Original objects regarding the reaction, when called on a `Reply`.
	OriginalReply      *Reply
	OriginalAckMessage *slack.AckMessage

	// When called on `Message`
	OriginalMessage *Message

	// Listener is a reference to the thing listening for incoming Reactions
	// you can call .Close() on it after a certain amount of time or after
	// the user you were interested in processed its things.
	Listener *ReactionListener
}

type reaction int

// ReactionAdded is used as the `Type` field of `ReactionEvent` (which
// you can register with `Reply.OnReaction()`)
const ReactionAdded = reaction(2)

// ReactionRemoved is the flipside of `ReactionAdded`.
const ReactionRemoved = reaction(1)

// ParseReactionEvent parses and normalizes reaction events to ReactionEvents
func ParseReactionEvent(event interface{}) *ReactionEvent {
	var re ReactionEvent
	switch ev := event.(type) {
	case *slack.ReactionAddedEvent:
		re.Type = ReactionAdded
		re.Emoji = ev.Reaction
		re.User = ev.User
		re.Item.Type = ev.Item.Type
		re.Item.Channel = ev.Item.Channel
		re.Item.File = ev.Item.File
		re.Item.FileComment = ev.Item.FileComment
		re.Item.Timestamp = ev.Item.Timestamp
		re.Timestamp = ev.EventTimestamp

	case *slack.ReactionRemovedEvent:
		re.Type = ReactionRemoved
		re.Emoji = ev.Reaction
		re.User = ev.User
		re.Item = ev.Item
		re.Item.Type = ev.Item.Type
		re.Item.Channel = ev.Item.Channel
		re.Item.File = ev.Item.File
		re.Item.FileComment = ev.Item.FileComment
		re.Item.Timestamp = ev.Item.Timestamp
		re.Timestamp = ev.EventTimestamp

	default:
		return nil
	}

	return &re
}
