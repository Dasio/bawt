---
title: "Types"
weight: 10
---

{{% children description="true" showhidden="true" depth="999" style="div" %}}

All plugin types are interfaces. Your plugin struct should satisfy the desired interface.

{{% notice tip %}}
You are not limited to using a singular interface.
{{% /notice %}}

## PluginInitializer

This is the basic chat plugin interface. The core of this interface is the `bawt.Listener` struct which provides extensive abilities to both listen to and respond to users in a variety of ways.

### Methods

| Method | Description |
| --- | :-- |
| `InitPlugin(*Bot)` | Used to load config and register listeners | 

## WebServer

The purpose of this plugin is to start a web server. Bawt comes with a built in web server with both a public and private router. The web server interface should not be used more than once.

### Methods

| Method | Description |
| --- | :-- |
| `InitWebServer(*Bot, []string)` | Used to load config and create routers |
| `RunServer()` | This method is called by Bawt to start up the web server in a go routine |
| `SetAuthMiddleware(func(http.Handler) http.Handler)` | Used to inject authentication middleware |
| `SetAuthenticatedUserFunc(func(req *http.Request) (*slack.User, error))` | The function to be executed when a user attempts to auth |
| `PrivateRouter() *mux.Router` | Returns a private router |
| `PublicRouter() *mux.Router` | Returns a public router |
| `GetSession(*http.Request) *sessions.Session` | Returns a session for an HTTP Request |
| `AuthenticatedUser(*http.Request) (*slack.User, error)` | Determines the authenticated user |

## WebPlugin

A web plugin is a plugin that can serve web pages on Bawt's web server.

### Methods

| Method | Description |
| --- | :-- |
| `InitWebPlugin(bot *Bot, private *mux.Router, public *mux.Router)` | Used to load config and register routes |

## WebServerAuth

WebServerAuth allows authentication to web services, such as through Slack, so that you can verify a user is who they say they are. WebServerAuth was mostly built on top of OAuth.

| Method | Description |
| --- | :-- |
| `InitWebServerAuth(bot *Bot, webserver WebServer)` | Initializes a web server auth plugin. Used to register the plugin and load config. |