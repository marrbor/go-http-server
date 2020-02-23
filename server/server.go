package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/marrbor/gohttp"
	"github.com/marrbor/golog"
)

var (
	ExitServerError = fmt.Errorf("unexpected exit ListenAndServe method")
)

// EntryPoint holds strings/function pair.
type EntryPoint struct {
	Resource string
	Function func(http.ResponseWriter, *http.Request)
}

// HttpServer holds configurator web server.
type HttpServer struct{ *http.Server }

// Start starts configurator web server. call as goroutine.
func (hs *HttpServer) Start(bus chan error) {
	if err := hs.ListenAndServe(); err != nil {
		bus <- err
		return
	}
	bus <- ExitServerError
}

// Stop stops web server.
func (hs *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := hs.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		golog.Error(err)
	}
	golog.Info("http server shutdown, exit.")
}

// NewServer returns new server instance.
func NewServer(port int, ep []EntryPoint, rf func(w http.ResponseWriter, r *http.Request)) *HttpServer {
	mux := http.NewServeMux()

	// register both "/entry/" and "/entry" points.
	for _, ep := range ep {
		mux.HandleFunc(fmt.Sprintf("/%s", ep.Resource), ep.Function)
		mux.HandleFunc(fmt.Sprintf("/%s/", ep.Resource), ep.Function)
	}

	if rf == nil {
		rf = func(w http.ResponseWriter, r *http.Request) { gohttp.NotFound(w, nil) }
	}
	mux.HandleFunc("/", rf)
	return &HttpServer{&http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}}
}
