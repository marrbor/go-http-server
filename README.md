# go-http-server

a lightweight http server.

`import "github.com/marrbor/go-http-server/server"`

## Usage

### Prepare an endpoint function.
Signature of endpoint function is `func xxx(w http.ResponseWriter, r *http.Request)`. no return values. 

e.g.:
```go
func foo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		gohttp.MethodNotAllowed(w, nil)
	}
	if err := gohttp.JSONResponse(w, &a); err != nil {
		gohttp.InternalServerError(w, err)
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	gohttp.responseOK()
}
```

### Prepare function entry point.
Its type is `EntryPoint` structure.

e.g.:
```go
var eps = []server.EntryPoint{
	{Resource: "foo", Function: foo},
	{Resource: "bar", Function: bar},
}

```

### Generate server instance and start it.

- Server takes one `chan error` channel at `Start`.
- When something wrong to run the server, the server send error to this channel.
- When `Stop` called, the server also return `http: Server closed` error.


```go
func main() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	hs := server.NewServer(8888, eps, nil)

	srvCh := make(chan error)
	hs.Start(srvCh)
LOOP:
    for {
		select {
			case sig := <- sigCh:
                golog.Info(sig)
				hs.Stop() // stop server when signal detect.
			case err := <- srvCh:
                golog.Info(err)
				break LOOP // exit loop when
		}
	}
	os.Exit(0)
}
```

# LICENSE
MIT
