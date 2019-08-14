package bawt

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nlopes/slack"
)

//
// Bot plugins
//

// Plugin describes the generic bot plugin
type Plugin interface{}

// Command is a command the bot is capable of understanding
type Command struct {
	Usage    string
	HelpText string
}

// PluginInitializer describes the interface is used to check which plugins
// can be initialized during plugin initalization initChatPlugins
type PluginInitializer interface {
	InitPlugin(*Bot)
}

// WebServer describes the interface for webserver plugins
type WebServer interface {
	// Used internally by the `bawt` library.
	InitWebServer(*Bot, []string)
	RunServer()

	// Used by an Auth provider.
	SetAuthMiddleware(func(http.Handler) http.Handler)
	SetAuthenticatedUserFunc(func(req *http.Request) (*slack.User, error))

	// Can be called by any plugins.
	PrivateRouter() *mux.Router
	PublicRouter() *mux.Router
	GetSession(*http.Request) *sessions.Session
	AuthenticatedUser(*http.Request) (*slack.User, error)
}

// WebPlugin initializes plugins with a `Bot` instance, a `privateRouter` and a `publicRouter`. All URLs handled by the `publicRouter` must start with `/public/`.
type WebPlugin interface {
	InitWebPlugin(bot *Bot, private *mux.Router, public *mux.Router)
}

// WebServerAuth returns a middleware warpping the passed on `http.Handler`. Only one auth handler can be added.
type WebServerAuth interface {
	InitWebServerAuth(bot *Bot, webserver WebServer)
}

var registeredPlugins = make([]Plugin, 0)

// RegisterPlugin adds the provided Plugin to the list of registered plugins
func RegisterPlugin(plugin Plugin) {
	registeredPlugins = append(registeredPlugins, plugin)
}

// RegisteredPlugins returns the list of registered plugins
func RegisteredPlugins() []Plugin {
	return registeredPlugins
}

func initPlugins(bot *Bot) {
	var enabledPlugins []string

	log := bot.Logging.Logger

	for _, plugin := range registeredPlugins {
		pluginType := reflect.TypeOf(plugin)
		if pluginType.Kind() == reflect.Ptr {
			pluginType = pluginType.Elem()
		}
		var typeList []string
		if _, ok := plugin.(PluginInitializer); ok {
			typeList = append(typeList, "Plugin")
		}
		if _, ok := plugin.(WebServer); ok {
			typeList = append(typeList, "WebServer")
		}
		if _, ok := plugin.(WebServerAuth); ok {
			typeList = append(typeList, "WebServerAuth")
		}
		if _, ok := plugin.(WebPlugin); ok {
			typeList = append(typeList, "WebPlugin")
		}

		log.Infof("Plugin %s implements %s", pluginType.String(),
			strings.Join(typeList, ", "))
		enabledPlugins = append(enabledPlugins, strings.Replace(pluginType.String(), ".", "_", -1))
	}

	initWebServer(bot, enabledPlugins)
	initWebPlugins(bot)

	if bot.WebServer != nil {
		go bot.WebServer.RunServer()
	}

	initChatPlugins(bot)
}

func initChatPlugins(bot *Bot) {
	for _, plugin := range registeredPlugins {
		chatPlugin, ok := plugin.(PluginInitializer)
		if ok {
			chatPlugin.InitPlugin(bot)
		}
	}
}

func initWebServer(bot *Bot, enabledPlugins []string) {
	for _, plugin := range registeredPlugins {
		webServer, ok := plugin.(WebServer)
		if ok {
			webServer.InitWebServer(bot, enabledPlugins)
			bot.WebServer = webServer
			return
		}
	}
}

func initWebPlugins(bot *Bot) {
	log := bot.Logging.Logger
	// If the WebServer isn't configured then don't load WebPlugins
	if bot.WebServer == nil {
		return
	}

	for _, plugin := range registeredPlugins {
		if webPlugin, ok := plugin.(WebPlugin); ok {
			webPlugin.InitWebPlugin(bot, bot.WebServer.PrivateRouter(), bot.WebServer.PublicRouter())
		}

		count := 0
		if webServerAuth, ok := plugin.(WebServerAuth); ok {
			count++

			if count > 1 {
				log.Fatal("Can not load two WebServerAuth plugins. Already loaded one.")
			}
			webServerAuth.InitWebServerAuth(bot, bot.WebServer)
		}
	}
}
