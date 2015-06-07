# einhorn
--
    import "github.com/dcu/http-einhorn"

Package einhorn contains a series of helpers to run a http server on stripe's
einhorn. It's helpful to run zero-downtime deployments.

For example to run `gin` with einhorn you can do the following:

    r := gin.Default()
    ...
    if einhorn.IsRunning() {
        einhorn.Start(r, 0)
    } else {
        r.Run(":8000")
    }

and then

    go build your-gin-app.go
    einhorn -b 0.0.0.0:4000 -- ./your-gin-app

please note you have to build the application otherwise the restart signal is
not handled. Now try restarting the cluster with:

    einhornsh -e "upgrade"


### Tests

To run the tests install gocheck:

    go get -u gopkg.in/check.v1

then run the tests as you would normally:

    go test .

## Usage

```go
var (
	// StopDelay indicates the time the process have to wait before stopping the process.
	StopDelay = 5 * time.Second

	// OnStopCallback is an optional callback that's called before closing the server.
	// Use it to cleanup everything before the server is stopped.
	OnStopCallback func(server *http.Server, listener net.Listener)
)
```

#### func  IsRunning

```go
func IsRunning() bool
```
IsRunning returns true if einrhorn is running this process.

#### func  Run

```go
func Run(httpServer *http.Server, einhornListenerFd int) error
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
