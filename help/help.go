package help

import (
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gopherworks/bawt"
)

// Help represents the help configuration
type Help struct {
	bot *bawt.Bot
}

type app struct {
	name        string
	description string
}

func init() {
	bawt.RegisterPlugin(&Help{})
}

// InitPlugin initializes the plugin
func (h *Help) InitPlugin(bot *bawt.Bot) {
	h.bot = bot

	h.listenHelp()
}

func (h *Help) listenHelp() {
	h.bot.Listen(&bawt.Listener{
		Matches:            regexp.MustCompile(`^!help.*`),
		MessageHandlerFunc: h.handleHelp,
		Name:               "Help",
		Description:        "Provides useful information about the apps and commands available",
		Commands: []bawt.Command{
			{
				Usage:    "!help",
				HelpText: "Displays the help topics for all registered plugins",
			},
			{
				Usage:    "!help <slug>",
				HelpText: "Displays the help topic for a particular plugin",
			},
		},
	})

	h.bot.Listen(&bawt.Listener{
		Matches:            regexp.MustCompile(`^!apps`),
		MessageHandlerFunc: h.handleApps,
		Name:               "Help",
		Description:        "Provides Information",
		Commands: []bawt.Command{
			{
				Usage:    "!apps",
				HelpText: "Displays a list of plugins",
			},
		},
	})

	h.bot.Listen(&bawt.Listener{
		Matches:            regexp.MustCompile(`^!bawt.*`),
		MessageHandlerFunc: h.handleBawt,
		Name:               "Help",
		Description:        "Provides Information",
		Commands: []bawt.Command{
			{
				Usage:    "!bawt",
				HelpText: "Displays a list of plugins",
			},
		},
	})
}

// It's important to remember that the global help is and always will be opt-in
func (h *Help) handleHelp(listen *bawt.Listener, msg *bawt.Message) {
	msg.AddReaction("+1") // Let the user know we're processing their request
	listeners := h.bot.Listeners()

	for _, l := range listeners {
		for _, c := range l.Commands {
			msg.Reply("%s\t\t%s", c.Usage, c.HelpText)
		}
	}
}

func (h *Help) handleApps(listen *bawt.Listener, msg *bawt.Message) {
	msg.AddReaction("+1")
	listeners := h.bot.Listeners()

	apps := []app{}

	// Deduplicating the app list
	for _, l := range listeners {
		d := false
		e := app{
			name:        l.Name,
			description: l.Description,
		}

		for _, a := range apps {
			if e.name == a.name {
				d = true
			}
		}

		if !d {
			apps = append(apps, e)
		}
	}

	for _, a := range apps {
		msg.Reply("%s\t\t%s", a.name, a.description)
	}
}

func (h *Help) handleBawt(listen *bawt.Listener, msg *bawt.Message) {
	msg.AddReaction("robot_face")

	parts := strings.Split(msg.Match[0], " ")

	act := parts[1]

	switch act {
	case "version":
		msg.Reply("*bawt* `v%s` (https://github.com/gopherworks/bawt/releases/tag/v%s)", bawt.Version, bawt.Version)
	case "dump-config":
		c := spew.Sprintf("%#+v", h.bot)
		p := bawt.FileUploadParameters{
			Content:        c,
			Filetype:       "Go",
			Filename:       "bot.go",
			Title:          "bawt.Bot{}",
			InitialComment: "This is a live snapshot of my config. This may contain sensitive data.",
		}
		p.Channels = append(p.Channels, msg.FromChannel.ID)

		msg.ReplyWithFile(p)
	}
}
