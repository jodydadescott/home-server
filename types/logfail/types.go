package logfail

import (
	"strings"
)

// LogFail should we log or should we crash
type LogFail string

const (
	Empty   LogFail = ""
	Log             = "log"
	Fail            = "Fail"
	Invalid         = "INVALID"
)

// NewFromString returns enum value from string
func NewFromString(input string) LogFail {

	switch strings.ToLower(input) {

	case string(Log):
		return Log

	case string(Fail):
		return Fail

	case "":
		return Empty

	}

	return Invalid
}
