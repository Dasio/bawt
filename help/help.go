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
		FromInternalGroup:  []string{"GlobalAdmins"},
		Commands: []bawt.Command{
			{
				Usage:    "!bawt",
				HelpText: "Displays a list of plugins",
			},
			{
				Usage:    "!bawt group list",
				HelpText: "Displays a list of groups",
			},
			{
				Usage:    "!bawt group <group> add-user",
				HelpText: "Add a user to a group",
			},
			{
				Usage:    "!bawt group <group> remove-user",
				HelpText: "Remove a user from a group",
			},
			{
				Usage:    "!bawt group <group> list-users",
				HelpText: "List the users in a group",
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
			// Found a duplicate
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
	const action = 1
	const user = 2

	log := h.bot.Logging.Logger

	msg.AddReaction("robot_face")

	parts := strings.Split(msg.Match[0], " ")

	if len(parts) == 1 {
		msg.Reply("Looks like you're missing an argument! Maybe consider `!help`?")
		return
	}

	a := parts[action]

	switch a {
	case "version":
		msg.Reply("*bawt* `v%s` (Release Notes: https://github.com/gopherworks/bawt/releases/tag/v%s)", bawt.Version, bawt.Version)
	case "dump-config":
		s := spew.ConfigState{
			Indent: "\t",
		}

		c := s.Sprintf("%#v", h.bot)
		p := bawt.FileUploadParameters{
			Content:        c,
			Filetype:       "Go",
			Filename:       "bot.go",
			Title:          "bawt.Bot{}",
			InitialComment: "This is a live snapshot of my config. This may contain sensitive data.",
		}
		p.Channels = append(p.Channels, msg.FromChannel.ID)

		msg.ReplyWithFile(p)
	case "whois":

		u := parts[user]
		u = bawt.NormalizeID(u)

		usr, err := h.bot.Slack.GetUserInfo(u)
		if err != nil {
			// We've reached an error
			log.WithError(err).Errorf("Error retrieving user info for %s", u)
			msg.Reply("User not found")

			return
		}

		// We found the user
		msg.Reply("Their user ID is %s", usr.ID)
	case "whoami":
		u := msg.FromUser

		msg.Reply("Your real name is %s (User: %s/ID: %s). You live in the %s timezone. Admin: %b; Owner: %b; Primary Owner: %b", u.RealName, u.Name, u.ID, u.TZLabel, u.IsAdmin, u.IsOwner, u.IsPrimaryOwner)
	case "channels":
		chans := []string{}

		for _, c := range h.bot.Channels {
			if c.IsChannel {
				chans = append(chans, c.Name)
			}
		}

		msg.Reply("I'm in the following channels: %s", strings.Join(chans, ", "))
	case "group":
		h.handleGroup(listen, msg)
	default:
		msg.Reply("I didn't recognize that argument! Maybe consider `!help`?")
	}
}

func (h *Help) handleGroup(listen *bawt.Listener, msg *bawt.Message) {
	const group = 2
	const action = 3
	const user = 4

	log := h.bot.Logging.Logger

	parts := strings.Split(msg.Match[0], " ")

	if len(parts) == 2 {
		msg.Reply("Looks like you're missing an argument! Maybe consider `!help`?")
		return
	}

	g := parts[group]
	a := parts[action]
	u := bawt.NormalizeID(parts[user])

	switch a {
	case "add-user":
		g := bawt.InternalGroup{
			Name: g,
		}

		g.Get(h.bot.DB)

		member, err := g.IsUserMember(h.bot.DB, msg.FromUser.ID)
		if err != nil {
			log.WithError(err).Error("Error determing user membership")
			return
		}

		if !member {
			msg.Reply("You don't have the proper permissions to do that.")
			return
		}

		if g.Name == "GlobalAdmins" {
			msg.Reply("GlobalAdmins cannot be modified via chat.")
			return
		}

		if g.FindDuplicate(h.bot.DB, u) {
			msg.Reply("That user is already a member of that group.")
			return
		}

		g.AddMember(h.bot.DB, u)
	case "remove-user":
		g := bawt.InternalGroup{
			Name: g,
		}

		g.Get(h.bot.DB)

		member, err := g.IsUserMember(h.bot.DB, msg.FromUser.ID)
		if err != nil {
			return
		}

		if !member {
			msg.Reply("You don't have the proper permissions to do that.")
			return
		}

		if g.Name == "GlobalAdmins" {
			msg.Reply("GlobalAdmins cannot be modified via chat.")
			return
		}

		if !g.FindDuplicate(h.bot.DB, u) {
			msg.Reply("That user is not a member of that group.")
			return
		}

		if u == msg.FromUser.ID {
			msg.Reply("You cannot remove yourself from a group.")
			return
		}

		g.RemoveMember(h.bot.DB, u)
	default:
		msg.Reply("I didn't understand your message.")
	}

}
