package bawt

// Config holds the configuration to connect with a given slack organization
type Config struct {
	JoinChannels   []string `json:"join_channels" mapstructure:"join_channels"`
	GeneralChannel string   `json:"general_channel" mapstructure:"general_channel"`
	TeamDomain     string   `json:"team_domain" mapstructure:"team_domain"`
	TeamID         string   `json:"team_id" mapstructure:"team_id"`
	APIToken       string   `json:"api_token" mapstructure:"api_token"`
	WebBaseURL     string   `json:"web_base_url" mapstructure:"web_base_url"`
	DBPath         string   `json:"db_path" mapstructure:"db_path"`
	PIDPath        string   `json:"pid_path" mapstructure:"pid_path"`
}
