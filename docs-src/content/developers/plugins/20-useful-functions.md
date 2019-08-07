---
title: "Loading Config"
weight: 30
---

Bawt lets plugins consume configuration just the same way Bawt does.

`bawt.LoadConfig()` can be utilized to marshal a struct that has fields you're looking for. Let's take a look at what that looks like in practice:

Let's say we're developing a plugin called Sample Plugin

```go
type SamplePlugin struct {
    bot     *bawt.Bot
}
```

SamplePlugin needs some config though.

```go
type PluginConfig struct {
	Foo      string `json:"foo"`
	Bar      string `json:"bar"`
	Foobar   string `json:"foobar"`
}
```

Bawt utilizes Viper under the hood, so if you're familiar with Viper then you'll find working with Bawt's config loader very simple. Bawt can accept JSON, YAML, or TOML.

This config needs to be stored somewhere! Let's modify our plugin struct:

```go
type SamplePlugin struct {
    bot       *bawt.Bot
    config    *PluginConfig
}
```

The InitPlugin method is called directly after the plugin has loaded. This is a good place to wire up our configuration:

```go
func (s *SamplePlugin) InitPlugin(b *bawt.Bot) {
    // The field name here is how our config field will be named
    var conf struct {
		SamplePlugin PluginConfig
	}
	
    // Drop bot into the plugin struct
    s.bot = b

    bot.LoadConfig(&conf)

    s.config = &conf.SamplePlugin
}
```

{{% notice tip %}}
This method utilizes a pointer receiver because it modifies our plugin struct.
{{% /notice %}}

Now if our config goes like this...

```yaml
Config:
    ...
SamplePlugin:
    foo: bar
    bar: foo
    foobar: barfoo
```

Bawt will pick up the SamplePlugin object and unmarshal it to your struct which is accessible by

```go
func (s SamplePlugin) someFunc() {
    fmt.Printf("Sample plugin foo: %s", s.config.foo)
}
```

You've now successfully made Bawt load plugin config on your behalf!