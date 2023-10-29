package proto

import (
	"strings"
)

// Proto is the protocol type. Currently only UDP and TCP are supported.
type Proto string

const (
	Empty   Proto = ""
	UDP           = "udp"
	TCP           = "tcp"
	Invalid       = "INVALID"
)

// NewFromString returns enum value from string
func NewFromString(input string) Proto {

	switch strings.ToLower(input) {

	case string(UDP):
		return UDP

	case string(TCP):
		return TCP

	case "":
		return Empty

	}

	return Invalid
}
