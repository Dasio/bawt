package bawt

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/nlopes/slack"
)

// Groups is the name of the Bolt DB bucket for groups
const Groups = "groups"

// GroupMembers is the name of the Bolt DB key for members of groups
const GroupMembers = "Members"

// GroupEmptyList is used to declare an empty list
const GroupEmptyList = "[\"\"]"

// GroupEmptyObject is used to declare an empty JSON object
const GroupEmptyObject = "{}"

// GroupSlackGroup is the name of the Bolt DB key for an Internal Groups corresponding Slack Group (if any)
const GroupSlackGroup = "SlackGroup"

// InternalGroup represents a group internal to the framework
type InternalGroup struct {
	Name       string
	SlackGroup slack.UserGroup `json:",omitempty"`
	Members    []string
}

// IsUserMember looks for a user ID that is a member of the given group
func (g InternalGroup) IsUserMember(db *bolt.DB, user string) (bool, error) {
	if err := g.Get(db); err != nil {
		return false, err
	}

	for _, m := range g.Members {
		if user == m {
			return true, nil
		}
	}

	return false, nil
}

// AddMember appends a user to the member list
func (g *InternalGroup) AddMember(db *bolt.DB, user string) error {
	if err := g.Get(db); err != nil {
		return err
	}

	if g.FindDuplicate(db, user) {
		return nil
	}

	g.Members = append(g.Members, user)

	if err := g.Put(db); err != nil {
		return err
	}

	return nil
}

// FindDuplicate returns true if it finds a duplicate
func (g InternalGroup) FindDuplicate(db *bolt.DB, user string) bool {
	for _, u := range g.Members {
		if user == u {
			return true
		}
	}

	return false
}

// RemoveMember removes a user from the members list
func (g *InternalGroup) RemoveMember(db *bolt.DB, user string) error {
	if err := g.Get(db); err != nil {
		return err
	}

	for i, m := range g.Members {
		if user == m {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
		}
	}

	if err := g.Put(db); err != nil {
		return err
	}

	return nil
}

// Get fetches the data from the database and unmarshals it into the struct
func (g *InternalGroup) Get(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// This will always exist because it's created if it doesn't exist at runtime
		b := tx.Bucket([]byte(Groups))

		// Ensuring that when plugins call a group that doesn't exist we at least instantiate it
		grp, err := b.CreateBucketIfNotExists([]byte(g.Name))
		if err != nil {
			return err
		}

		// Fetch and store members
		m := grp.Get([]byte(GroupMembers))

		// If the group doesn't exist then set the default Members
		if len(m) == 0 {
			m = []byte(GroupEmptyList)
		}

		err = json.Unmarshal(m, &g.Members)
		if err != nil {
			return err
		}

		// Fetch and store the slack group
		m = grp.Get([]byte(GroupSlackGroup))

		// If the group doesn't exist then set the default Members
		if len(m) == 0 {
			m = []byte(GroupEmptyObject)
		}

		err = json.Unmarshal(m, &g.SlackGroup)
		if err != nil {
			return err
		}

		return nil
	})
}

// Put pulls information out of the struct and stores it in the database
func (g *InternalGroup) Put(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// This will always exist because it's created if it doesn't exist at runtime
		b := tx.Bucket([]byte(Groups))

		// Ensuring that when plugins call a group that doesn't exist we at least instantiate it
		grp, err := b.CreateBucketIfNotExists([]byte(g.Name))
		if err != nil {
			return err
		}

		// Fetch and store members
		m, err := json.Marshal(g.Members)
		if err != nil {
			return err
		}

		// If the group doesn't exist then set the default Members
		if len(m) == 0 {
			m = []byte(GroupEmptyList)
		}

		grp.Put([]byte(GroupMembers), m)

		// Fetch and store the slack group
		m, err = json.Marshal(g.SlackGroup)
		if err != nil {
			return err
		}

		// If the group doesn't exist then set the default Members
		if len(m) == 0 {
			m = []byte(GroupEmptyList)
		}

		grp.Put([]byte(GroupSlackGroup), m)

		return nil
	})
}
