package server

import (
	"context"

	logger "github.com/jodydadescott/jody-go-logger"
	"go.uber.org/zap"

	"github.com/jodydadescott/home-server/dns"
	"github.com/jodydadescott/home-server/http"
	"github.com/jodydadescott/home-server/static"
	"github.com/jodydadescott/home-server/types"
	"github.com/jodydadescott/home-server/unifi"
)

type Config = types.Config

type Server struct {
	dns  *dns.Server
	http *http.Server
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

	s := &Server{
		dns: dns.New(dnsConfig),
	}

	if config.HttpConfig != nil && config.HttpConfig.Enabled {
		zap.L().Debug("HTTP Server is enabled")

		httpConfig := &http.Config{
			Listener:       config.HttpConfig.Listener,
			RecordProvider: s.dns,
		}

		s.http = http.New(httpConfig)

	} else {
		zap.L().Debug("HTTP Server is not enabled")
	}

	return s
}

func (t *Server) Run(ctx context.Context) error {

	errs := make(chan error, 2)

	dnsCtx, dnsCancel := context.WithCancel(ctx)
	httpCtx, httpCancel := context.WithCancel(ctx)

	defer func() {
		dnsCancel()
		httpCancel()
	}()

	go func() {
		err := t.dns.Run(dnsCtx)
		if err != nil {
			errs <- err
		}
	}()

	if t.http != nil {
		go func() {
			err := t.http.Run(httpCtx)
			if err != nil {
				errs <- err
			}
		}()
	}

	select {

	case err := <-errs:
		zap.L().Info("Shutting down or error")
		return err

	case <-ctx.Done():
		zap.L().Info("Shutting down or signal")

	}

	return nil
}
