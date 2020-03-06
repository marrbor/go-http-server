// simple http server that only shows received request information.
// set listen port via argument or environment variable "GO_HTTP_SERVER_PORT".
// when specified port number by both, use argument prior to environment variable.
// the default value of port is 8765.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/marrbor/go-http-server/server"
	"github.com/marrbor/gohttp"
	"github.com/marrbor/golog"
)

const (
	DefaultPort = 8765
	EnvPort     = "GO_HTTP_SERVER_PORT"
)

func all(w http.ResponseWriter, r *http.Request) {
	golog.Info(fmt.Sprintf("Recevied from: %s", r.RemoteAddr))
	golog.Info(fmt.Sprintf("Method: %s", r.Method))

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gohttp.BadRequest(w, err)
		return
	}
	golog.Info(fmt.Sprintf("Body: %s", string(b)))
	gohttp.ResponseOK(w)
}

var eps = []server.EntryPoint{
	{Resource: "/", Function: all},
}

func main() {
	port := DefaultPort
	if 0 < len(os.Getenv(EnvPort)) {
		p, err := strconv.Atoi(os.Getenv(EnvPort))
		if err != nil {
			golog.Panic(err)
		}
		port = p
	}

	if 1 < len(os.Args) {
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			golog.Panic(err)
		}
		port = p
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	hs := server.NewServer(port, nil, all)
	srvCh := make(chan error)
	hs.Start(srvCh)
LOOP:
	for {
		select {
		case sig := <-sigCh:
			golog.Info(sig)
			hs.Stop()
		case err := <-srvCh:
			golog.Info(err)
			break LOOP
		}
	}
	os.Exit(0)
}
