package bawt

import (
	"math/rand"
	"time"
)

var registeredStrings = make(map[string][]string)

// RegisterStringList takes an array of strings and stores them in a category
func RegisterStringList(category string, list []string) {
	registeredStrings[category] = list
}

// RandomString returns a random string from the array stored at category
func RandomString(category string) string {
	strList, ok := registeredStrings[category]
	if !ok {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	idx := rand.Int() % len(strList)
	return strList[idx]
}
