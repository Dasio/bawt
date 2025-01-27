package bawt

import (
	"time"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

// ReplyWithFile replies to a user with a file
type ReplyWithFile struct {
	*slack.File
	bot *Bot
}

// Reply represents a reply to bawt
type Reply struct {
	*slack.OutgoingMessage
	bot *Bot
}

// AddReaction adds a reaction to a reply
func (r *Reply) AddReaction(emoji string) *Reply {
	r.OnAck(func(ev *slack.AckMessage) {
		go r.bot.Slack.AddReaction(emoji, slack.NewRefToMessage(r.Channel, ev.Timestamp))
	})
	return r
}

// DeleteAfter deletes a reply after a certain duration
func (r *Reply) DeleteAfter(duration string) *Reply {
	timeDur := parseAutodestructDuration("DeleteAfter", duration, r.bot.Logging.Logger)

	r.OnAck(func(ev *slack.AckMessage) {
		go func() {
			time.Sleep(timeDur)
			r.bot.Slack.DeleteMessage(r.Channel, ev.Timestamp)
		}()
	})

	return r
}

// ListenReaction listens for reactions
func (r *Reply) ListenReaction(reactListen *ReactionListener) {
	r.OnAck(func(ackEv *slack.AckMessage) {
		listen := reactListen.newListener()
		listen.EventHandlerFunc = func(_ *Listener, event interface{}) {
			re := ParseReactionEvent(event)
			if re == nil {
				return
			}

			if ackEv.Timestamp != re.Item.Timestamp {
				return
			}

			if re.User == r.bot.Myself.ID {
				return
			}

			if !reactListen.filterReaction(re) {
				return
			}

			re.OriginalReply = r
			re.OriginalAckMessage = ackEv
			re.Listener = reactListen

			reactListen.HandlerFunc(reactListen, re)
		}
		r.bot.Listen(listen)
	})
}

// OnAck allows you to catch the message_id of the message you
// replied.  Call it immediately after sending a reply to be sure to
// catch the confirmation message with the message_id.
//
// With the message_id, you can modify your reply, add reactions to it
// or delete it.
func (r *Reply) OnAck(f func(ack *slack.AckMessage)) {
	log := r.bot.Logging.Logger

	r.bot.Listen(&Listener{
		ListenDuration: 20 * time.Second,
		EventHandlerFunc: func(subListen *Listener, event interface{}) {
			if ev, ok := event.(*slack.AckMessage); ok {
				if ev.ReplyTo == r.ID {
					f(ev)
					subListen.Close()
				}
			}
		},
		TimeoutFunc: func(subListen *Listener) {
			log.Println("OnAck Listener dropped, because no corresponding AckMessage was received before timeout")
			subListen.Close()
		},
	})
}

// Updateable returns an instance of UpdateableReply, which has a few
// methods to update a message after the fact.  It is safe to use in
// different goroutines no matter when.
func (r *Reply) Updateable() *UpdateableReply {
	updt := &UpdateableReply{
		reply: r,
	}

	r.OnAck(func(ack *slack.AckMessage) {
		updt.lock.Lock()
		defer updt.lock.Unlock()

		updt.msgTimestamp = ack.Timestamp
		go updt.dispatch()
	})

	return updt
}

// Listen here on Reply is the same as Bot.Listen except that
// ReplyAck() will be filled with the slack.AckMessage before any
// event is dispatched to this listener.
func (r *Reply) Listen(listen *Listener) error {
	log := r.bot.Logging.Logger

	listen.Bot = r.bot

	err := listen.checkParams()
	if err != nil {
		log.Println("Reply.Listen(): Invalid Listener: ", err)
		return err
	}

	r.OnAck(func(ev *slack.AckMessage) {
		listen.replyAck = ev
		r.bot.addListener(listen)
	})

	return nil
}

func parseAutodestructDuration(funcName string, duration string, log *logrus.Logger) time.Duration {
	timeDur, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("error: %s called with invalid `duration`: %q, using 1 second instead.\n", funcName, duration)
		timeDur = 1 * time.Second
	}
	return timeDur
}
