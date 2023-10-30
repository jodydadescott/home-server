package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"

	"github.com/jodydadescott/home-server/types"
)

type NetPort = types.NetPort

type HTTPRequest struct {
	Host   string      `json:"host,omitempty"`
	Method string      `json:"method,omitempty"`
	URL    *url.URL    `json:"url,omitempty"`
	Header http.Header `json:"header,omitempty"`
}

type Server struct {
	s              *http.Server
	recordProvider RecordProvider
}

// NewServer ...
func New(config *Config) *Server {

	if config == nil {
		panic("config is nil")
	}

	if config.Listener == nil {
		panic("NetPort is nil")
	}

	if config.RecordProvider == nil {
		panic("RecordProvider is required")
	}

	s := &Server{
		recordProvider: config.RecordProvider,
	}
	s.s = &http.Server{Addr: config.Listener.GetIPColonPort(), Handler: s}
	return s
}

func (t *Server) Run(ctx context.Context) error {

	go func() {
		<-ctx.Done()
		zap.L().Info("Shutting down HTTP server on signal")
		t.s.Shutdown(ctx)
	}()

	zap.L().Info("Starting HTTP Server")
	return t.s.ListenAndServe()
}

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	filter := r.URL.Query().Get("filter")

	switch r.URL.Path {

	case "/getdevices":
		w.Header().Set("Content-Type", "application/json")

		write := func(records *DomainRecords) {
			j, err := json.Marshal(records)
			if err != nil {
				zap.L().Error(err.Error())
			}
			fmt.Fprintf(w, string(j))
		}

		records := t.recordProvider.GetRecords()

		if filter == "" {
			write(records)
			return
		}

		newRecords := &DomainRecords{}

		for _, record := range records.ARecords {
			if strings.HasPrefix(record.Hostname, filter) {
				newRecords.AddARecords(record)
			}
		}

		for _, record := range records.AAAARecords {
			if strings.HasPrefix(record.Hostname, filter) {
				newRecords.AddAAAARecords(record)
			}
		}

		write(newRecords)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "<p>Hello</p>")
	fmt.Fprintf(w, "<p>You probably want to make one of the following calls</p>")
	fmt.Fprintf(w, fmt.Sprintf("<p><a href=\"http:/%s/getdevices\">/getdevices?filter=shelly</a></p>", r.Host))

}
