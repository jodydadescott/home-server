package types

import (
	"strings"
)

func cleanHostname(input string) string {
	hostname := strings.ToLower(input)
	hostname = strings.ToLower(hostname)
	hostname = space.ReplaceAllString(hostname, "-")
	return hostname
}
