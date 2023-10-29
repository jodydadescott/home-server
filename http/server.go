package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// HTTPRequest ...
type HTTPRequest struct {
	Host   string      `json:"host,omitempty"`
	Method string      `json:"method,omitempty"`
	URL    *url.URL    `json:"url,omitempty"`
	Header http.Header `json:"header,omitempty"`
}

// HTTPServer ...
type HTTPServer struct {
	s *http.Server
}

// NewServer ...
func NewServer() *HTTPServer {
	zap.L().Debug("NewServer()")
	s := &HTTPServer{}
	return s
}

// Listen ...
func (t *HTTPServer) Listen(listen string) error {
	zap.L().Debug("Listen() / Blocking")

	if t.s != nil {
		return fmt.Errorf("Duplicate call to Listen")
	}

	t.s = &http.Server{Addr: listen, Handler: t}

	err := t.s.ListenAndServe()
	if err != nil {
		zap.L().Debug("Listen() no longer blocking and returned error")
		return err
	}

	zap.L().Debug("Listen() no longer blocking")
	return nil
}

// Shutdown ...
func (t *HTTPServer) Shutdown() {
	zap.L().Debug("Shutdown()")

	if t.s == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t.s.Shutdown(ctx)
}

func pingResponse(r *http.Request) string {

	httpRequest := &HTTPRequest{
		Host:   r.Host,
		Method: r.Method,
		URL:    r.URL,
		Header: r.Header,
	}

	b, _ := json.Marshal(httpRequest)
	return string(b)
}

// ServeHTTP ...
func (t *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	j := pingResponse(r)
	zap.L().Debug(fmt.Sprintf("ServeHTTP(http.ResponseWriter, *http.Request=%s)", j))

	switch r.URL.Path {

	case "/":
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html>")
		fmt.Fprintf(w, "<h3>Welcome to the test server</h3>")
		fmt.Fprintf(w, "<a>All operations (GET, POST, PATCH, HEAD, PUT) are supported</p>")
		fmt.Fprintf(w, "<a href=\"http://"+r.Host+"/\">"+"http://"+r.Host+"/</a> is what you are looking at now</p>")
		fmt.Fprintf(w, "<a href=\"http://"+r.Host+"/info\">"+"http://"+r.Host+"/info</a> will return a JSON object with details about the request</p>")
		fmt.Fprintf(w, "<a>All other paths will print the path and operation")
		fmt.Fprintf(w, "</html>")
		return

	case "/info":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, j)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<a>Path: "+r.URL.Path+"</p>")
	fmt.Fprintf(w, "<a>Operation/Method: "+r.Method+"</p>")
	fmt.Fprintf(w, "</html>")
}
