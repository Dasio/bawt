package recognition

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gopherworks/bawt"
)

func (p *Plugin) listenUpvotes() {
	p.bot.Listen(&bawt.Listener{
		Name:        "Recognition",
		Description: "A fun game that encourages team mates to learn each others faces",
		EventHandlerFunc: func(_ *bawt.Listener, event interface{}) {
			react := bawt.ParseReactionEvent(event)
			if react == nil {
				return
			}

			log.Println("Fetching item ts:", react.Item.Timestamp)
			recognition := p.store.Get(react.Item.Timestamp)
			if recognition == nil {
				return
			}

			user := p.bot.Users[react.User]
			if user.IsBot {
				log.Println("Not taking votes from bots")
				return
			}

			if p.config.DomainRestriction != "" && !strings.HasSuffix(user.Profile.Email, p.config.DomainRestriction) {
				log.Printf("Not taking votes from people outsite domain %q, was %q", p.config.DomainRestriction, user.Profile.Email)
				return
			}

			log.Println("Up/down voting recognition")
			p.upvoteRecognition(recognition, react)
		},
	})
}

func (p *Plugin) upvoteRecognition(recognition *Recognition, reaction *bawt.ReactionEvent) {
	direction := 1
	if reaction.Type == bawt.ReactionRemoved {
		direction = -1
	}
	recognition.Reactions[reaction.User] += direction
	p.store.Put(recognition)
}
