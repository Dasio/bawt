// Package standup is a plugin for bawt that facilitates standups for teams
package standup

import "github.com/gopherworks/bawt"

type Standup struct {
	bot            *bawt.Bot
	sectionUpdates chan sectionUpdate
}

const TODAY = 0
const WEEKAGO = -6 // [0,-6] == 7 days

func init() {
	bawt.RegisterPlugin(&Standup{})
}

func (standup *Standup) InitPlugin(bot *bawt.Bot) {
	standup.bot = bot
	standup.sectionUpdates = make(chan sectionUpdate, 15)

	go standup.manageUpdatesInteraction()

	bot.Listen(&bawt.Listener{
		MessageHandlerFunc: standup.ChatHandler,
		Name:               "Standup",
		Description:        "Provides an assistant for running stand up over chat",
	})
}

func (standup *Standup) ChatHandler(listen *bawt.Listener, msg *bawt.Message) {
	res := sectionRegexp.FindAllStringSubmatchIndex(msg.Text, -1)
	if res != nil {
		for _, section := range extractSectionAndText(msg.Text, res) {
			standup.TriggerReminders(msg, section.name)
			// err := standup.StoreLine(msg, section.name, section.text)
			// if err != nil {
			// 	log.Println(err)
			// }
		}
	}
}
