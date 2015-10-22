# einhorn
--
    import "github.com/dcu/http-einhorn"

Package einhorn contains a series of helpers to run a http server on stripe's
einhorn. It's helpful to run zero-downtime deployments.

For example to run a web app with einhorn you can do the following:

```go
mux := http.NewServeMux()
mux.HandleFunc("/", httpHandler)
...
if einhorn.IsRunning() {
    einhorn.Start(mux, 0)
} else {
    server := &http.Server{Handler: mux, Addr: ":4000"}
     server.ListenAndServe()
}
```

and then

    go build your-app.go
    einhorn -b 0.0.0.0:4000 -- ./your-app

please note you have to build the application otherwise the restart signal is
not handled. Now try restarting the cluster with:

    einhornsh -e "upgrade"


### Tests

To run the tests install gocheck:

    go get -u gopkg.in/check.v1

then run the tests as you would normally:

    go test .


Integrating with Graceful

Graceful is a package enabling graceful shutdown of http.Handler servers.
Integration with graceful is easy, first create an instance of the graceful
server:

    gracefulServer := &graceful.Server{
        Server: http.Server{Handler: mux},
    }

then run it with einhorn:

    einhorn.Run(gracefulServer, 0)

also, you'll need to set the `OnStopCallback` to close the server:

    einhorn.OnStopCallback = func(server einhorn.Server, listener net.Listener) {
        gracefulServer := server.(*graceful.Server)
        gracefulServer.Stop(5 * time.Second)
    }

## Usage

```go
var (
	// StopDelay indicates the time the process have to wait before stopping the process.
	StopDelay = 5 * time.Second

	// OnStopCallback is an optional callback that's called before closing the server.
	// Use it to cleanup everything before the server is stopped.
	OnStopCallback func(server Server, listener net.Listener)
)
```

#### func  IsRunning

```go
func IsRunning() bool
```
IsRunning returns true if einrhorn is running this process.

#### func  Run

```go
func Run(httpServer Server, einhornListenerFd int) error
```
Run runs the given http server on the given `einhornListenerFd` The listener FD
is related to the einhorn's `-b` option so if you only pass one address the
einhornListenerFd should be 0. It also handles the restart signal.

#### func  Start

```go
func Start(handler http.Handler, einhornListenerFd int) error
```
Start starts the http handler on the given `einhornListenerFd`. The listener FD
is related to the einhorn's `-b` option so if you only pass one address the
einhornListenerFd should be 0.

#### type Server

```go
type Server interface {
	Serve(listener net.Listener) error
}
```

Server is a interface that requires to implement the Serve method. This
interface is compatible with http.Server and you can use it to make this package
compatible with other servers.
