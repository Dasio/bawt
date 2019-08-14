package bawt

import (
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

func TestListenerCheckParams(t *testing.T) {
	c := Listener{
		ListenUntil:    time.Now(),
		ListenDuration: 120 * time.Second,
	}
	err := c.checkParams()
	if err == nil {
		t.Error("checkParams shouldn't be nil")
	}
}

func TestDefaultFilter(t *testing.T) {
	b := &Bot{}
	b.Logging.Logger = logrus.New()
	b.Logging.Logger.Out = ioutil.Discard
	c := &Listener{}
	c.Bot = b
	u := &slack.User{ID: "a_user"}
	m := &Message{Msg: &slack.Msg{Text: "hello mama"}, FromUser: u}

	if c.filterMessage(m) != true {
		t.Error("filterMessage Failed")
	}

	type El struct {
		c *Listener
		r bool
	}
	tests := []El{
		{&Listener{Bot: b}, true},

		{&Listener{
			Bot:      b,
			Contains: "moo",
		}, false},

		{&Listener{
			Bot:      b,
			Contains: "MAMA",
		}, true},

		{&Listener{
			Bot:      b,
			FromUser: u,
		}, true},

		{&Listener{
			Bot:     b,
			Matches: regexp.MustCompile(`hello`),
		}, true},

		{&Listener{
			Bot:     b,
			Matches: regexp.MustCompile(`other-message`),
		}, false},

		{&Listener{
			Bot:      b,
			FromUser: &slack.User{ID: "another_user"},
		}, false},
	}

	for i, el := range tests {
		if el.c.filterMessage(m) != el.r {
			t.Error("filterMessage Failed, index ", i)
		}
	}
}

func TestMatchesMessage(t *testing.T) {
	c := &Listener{Bot: &Bot{Logging: Logging{Logger: logrus.New()}}, Matches: regexp.MustCompile(`(this) (is) (good)`)}
	c.Bot.Logging.Logger.Out = ioutil.Discard // Silence logs
	m := &Message{Msg: &slack.Msg{Text: "yeah this is good and all"}}

	if c.filterMessage(m) != true {
		t.Error("filterMessage Failed")
	}

	if len(m.Match) != 4 {
		t.Error("didn't find 4 matches")
	}

	if m.Match[0] != "this is good" {
		t.Error("didn't find 'this is good'")
	}

	if m.Match[1] != "this" {
		t.Error("didn't find 'this'")
	}
}
