package workerpool

import (
	"github.com/sapariduo/workpool/worker"
)

const (
	stateInitial = iota
	stateMain
	stateExit
	stateLaunch
	stateProcessing
	stateFinish
	stateSleep
	stateWait
	stateTimeout
	stateQuit
)

type transition struct {
	state   int
	payload payload
}

type payload struct {
	workerID     worker.ID
	isProcessing bool
}
