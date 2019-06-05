package bawt

import (
	"errors"
	"os"
)

// SlackConfig holds the configuration to connect with a given slack organization
type SlackConfig struct {
	Username       string
	Password       string
	Nickname       string
	JoinChannels   []string `json:"join_channels" mapstructure:"join_channels"`
	GeneralChannel string   `json:"general_channel" mapstructure:"general_channel"`
	TeamDomain     string   `json:"team_domain" mapstructure:"team_domain"`
	TeamID         string   `json:"team_id" mapstructure:"team_id"`
	APIToken       string   `json:"api_token" mapstructure:"api_token"`
	WebBaseURL     string   `json:"web_base_url" mapstructure:"web_base_url"`
	DBPath         string   `json:"db_path" mapstructure:"db_path"`
	Debug          bool
}

func checkPermission(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fi.Mode()&0077 != 0 {
		return errors.New("Config file is permitted to group/other. Do chmod 600 " + file)
	}
	return nil
}
