// Package einhorn contains a series of helpers to run a http server on stripe's einhorn.
// It's helpful to run zero-downtime deployments.
//
// For example to run `gin` with einhorn you can do the following:
//     r := gin.Default()
//     ...
//     if einhorn.IsRunning() {
//         einhorn.Start(r, 0)
//     } else {
//         r.Run(":8000")
//     }
// and then
//     go build your-gin-app.go
//     einhorn -b 0.0.0.0:4000 -- ./your-gin-app
// please note you have to build the application otherwise the restart signal is not handled.
// Now try restarting the cluster with:
//     einhornsh -e "upgrade"
//
// Tests
//
// To run the tests install gocheck:
//     go get -u gopkg.in/check.v1
// then run the tests as you would normally:
//     go test .
//
package einhorn

import (
	"github.com/stripe/go-einhorn/einhorn"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	// StopDelay indicates the time the process have to wait before stopping the process.
	StopDelay = 5 * time.Second

	// OnStopCallback is an optional callback that's called before closing the server.
	// Use it to cleanup everything before the server is stopped.
	OnStopCallback func(server *http.Server, listener net.Listener)
)

// IsRunning returns true if einrhorn is running this process.
func IsRunning() bool {
	einhornFDCount := os.Getenv("EINHORN_FD_COUNT")
	counter, err := strconv.Atoi(einhornFDCount)
	if err != nil {
		return false
	}
	return counter > 0
}

// Start starts the http handler on the given `einhornListenerFd`.
// The listener FD is related to the einhorn's `-b` option so if you
// only pass one address the einhornListenerFd should be 0.
func Start(handler http.Handler, einhornListenerFd int) error {
	httpServer := &http.Server{Handler: handler}

	return Run(httpServer, einhornListenerFd)
}

// Run runs the given http server on the given `einhornListenerFd`
// The listener FD is related to the einhorn's `-b` option so if you
// only pass one address the einhornListenerFd should be 0.
// It also handles the restart signal.
func Run(httpServer *http.Server, einhornListenerFd int) error {
	listener, err := einhorn.GetListener(einhornListenerFd)
	if err != nil {
		return err
	}
	defer listener.Close()

	restartChan := make(chan os.Signal, 1)
	signal.Notify(restartChan, syscall.SIGUSR2)
	go restartHandler(restartChan, httpServer, listener)

	err = httpServer.Serve(listener)
	opErr, ok := err.(*net.OpError)
	if ok && opErr.Op == "accept" {
		// ignore error produced when the listener is closed.
		return nil
	}

	return err
}

func restartHandler(restartChan chan os.Signal, server *http.Server, listener net.Listener) {
	<-restartChan
	server.SetKeepAlivesEnabled(false)

	if OnStopCallback != nil {
		OnStopCallback(server, listener)
	}
	time.Sleep(StopDelay)

	_ = listener.Close()
	signal.Stop(restartChan)

	close(restartChan)
}
