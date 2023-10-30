package types

import (
	"regexp"
	"time"

	"github.com/jodydadescott/home-server/types/proto"
)

const (
	DefaultDomain    = "home"
	DefaultDnsProto  = proto.UDP
	DefaultDnsPort   = 53
	DefaultDnsDomain = "home"
	DefaultRefresh   = time.Hour
	DefaultHTTPPort  = 8080
)

var space = regexp.MustCompile(`\s+`)
