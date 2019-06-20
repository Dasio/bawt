package bawt

// Mood is an enum
type Mood int

const (
	// Happy indicates a happy bot
	Happy Mood = iota
	// Hyper indicates a hyper bot
	Hyper
)

// WithMood returns a different response depending on the mood
func (bot *Bot) WithMood(happy, hyper string) string {
	if bot.Mood == Happy {
		return happy
	}

	return hyper
}
