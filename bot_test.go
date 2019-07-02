// Package bawt is a Slack bot framework written in Go
package bawt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// a is actual
	// e is expected

	// Test that the config location is stored
	a := New("/tmp/test.json")
	e := &Bot{configFile: "/tmp/test.json"}

	assert.Equal(t, e.configFile, a.configFile)
	assert.Empty(t, a.outgoingMsgCh)
	assert.Empty(t, a.outgoingFileCh)
	assert.Empty(t, a.addListenerCh)
	assert.Empty(t, a.delListenerCh)
	assert.Empty(t, a.Users)
	assert.Empty(t, a.Channels)
	assert.NotEmpty(t, a.PubSub)

	// Test storing an empty config
	a = New("")
	e = &Bot{configFile: ""}

	assert.Equal(t, e.configFile, a.configFile)
	assert.Empty(t, a.outgoingMsgCh)
	assert.Empty(t, a.outgoingFileCh)
	assert.Empty(t, a.addListenerCh)
	assert.Empty(t, a.delListenerCh)
	assert.Empty(t, a.Users)
	assert.Empty(t, a.Channels)
	assert.NotEmpty(t, a.PubSub)
}

func TestBot_readInConfig(t *testing.T) {
	bot := New("")

	e := "test"
	os.Setenv("CONFIG_API_TOKEN", e)
	bot.LoadConfig(bot)
	assert.NotEmpty(t, bot.Config.APIToken)
	assert.Equal(t, e, bot.Config.APIToken)

	e = "test"
	os.Setenv("CONFIG_GENERAL_CHANNEL", e)
	bot.LoadConfig(bot)
	assert.NotEmpty(t, bot.Config.GeneralChannel)
	assert.Equal(t, e, bot.Config.GeneralChannel)

	e = "test"
	os.Setenv("CONFIG_TEAM_DOMAIN", e)
	bot.LoadConfig(bot)
	assert.NotEmpty(t, bot.Config.TeamDomain)
	assert.Equal(t, e, bot.Config.TeamDomain)

	e = "test"
	os.Setenv("CONFIG_WEB_BASE_URL", e)
	bot.LoadConfig(bot)
	assert.NotEmpty(t, bot.Config.WebBaseURL)
	assert.Equal(t, e, bot.Config.WebBaseURL)

	e = "test"
	os.Setenv("CONFIG_DB_PATH", e)
	bot.LoadConfig(bot)
	assert.NotEmpty(t, bot.Config.DBPath)
	assert.Equal(t, e, bot.Config.DBPath)
}
