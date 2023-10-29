package server

import (
	"context"

	logger "github.com/jodydadescott/jody-go-logger"
	"go.uber.org/zap"

	"github.com/jodydadescott/home-server/dns"
	"github.com/jodydadescott/home-server/static"
	"github.com/jodydadescott/home-server/types"
	"github.com/jodydadescott/home-server/unifi"
)

type Config = types.Config

type Server struct {
	dns *dns.Server
}

func New(config *Config) *Server {

	trace := false

	if config.Logging != nil {
		logger.SetConfig(config.Logging)

		if config.Logging.LogLevel == logger.TraceLevel {
			trace = true
		}

	}

	dnsConfig := &dns.Config{
		Listeners:   config.Listeners,
		Nameservers: config.Nameservers,
		Trace:       trace,
	}

	if config.Unifi != nil && config.Unifi.Enabled {
		zap.L().Debug("Unifi is enabled")
		dnsConfig.AddProvider(unifi.New(config.Unifi))
	} else {
		zap.L().Debug("Unifi is not enabled")
	}

	if config.Static != nil && config.Static.Enabled {
		zap.L().Debug("static config is enabled")
		for _, v := range static.New(config.Static) {
			dnsConfig.AddProvider(v)
		}
	} else {
		zap.L().Debug("static config is not enabled")
	}

	return &Server{
		dns: dns.New(dnsConfig),
	}
}

func (t *Server) Run(ctx context.Context) error {
	return t.dns.Run(ctx)
}
