package workerpool

import (
	"os"
	"os/signal"
	"time"

	logger "github.com/uniplaces/go-logger"
)

type Shutdown struct {
	initiateChannel chan bool
	doneChannel     chan bool
	timeout         time.Duration
}

func NewShutdown(timeout time.Duration) *Shutdown {
	return &Shutdown{
		initiateChannel: make(chan bool, 1),
		doneChannel:     make(chan bool, 1),
		timeout:         timeout,
	}
}

func (shutdown *Shutdown) WaitForSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	logger.Info("received interrupt signal")
	shutdown.initiateChannel <- true
	select {
	case <-signalChannel:
		logger.Warning("forcing shutdown")
		os.Exit(1)
	case <-shutdown.doneChannel:
		logger.Info("cleanup successful, exiting")
	case <-time.After(time.Second * shutdown.timeout):
		logger.Info("cleanup timed out, exiting")
	}
}
