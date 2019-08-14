package bawt

import (
	"fmt"
	"strings"
)

// Status represents the status of various components
type Status struct {
	DB   string
	Chat string
	HTTP string
}

// NewStatus returns a new status struct
func NewStatus() Status {
	return Status{
		DB:   "Not ok",
		Chat: "Not ok",
		HTTP: "N/A",
	}
}

// Update updates the value of the component in the struct
func (s *Status) Update(comp string, value string) error {
	switch strings.ToLower(value) {
	case "ok":
		break
	case "not ok":
		break
	case "n/a":
		break
	default:
		return fmt.Errorf("Invalid value: %s", value)
	}

	switch strings.ToLower(comp) {
	case "db":
		s.DB = value
	case "chat":
		s.Chat = value
	case "http":
		s.HTTP = value
	default:
		return fmt.Errorf("Invalid component: %s", comp)
	}

	return nil
}
