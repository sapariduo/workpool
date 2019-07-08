package workerpool

import (
	"time"
	"github.com/sapariduo/workpool/worker"

	"github.com/uniplaces/carbon"
)

type WorkerPool struct {
	maxWorkers        int
	interval          int
	timeout           int
	shutdown          *Shutdown
	processingChannel chan bool
	finishedChannel   chan worker.ID
	needsToQuit       bool
	activeWorkers     map[worker.ID]carbon.Carbon
	currentTransition transition
}

func NewWorkerPool(config Config, shutdown *Shutdown) WorkerPool {
	maxWorkers := config.Workers
	initialTransition := transition{
		state: stateInitial,
	}

	return WorkerPool{
		maxWorkers:        maxWorkers,
		interval:          config.Interval,
		timeout:           config.Timeout,
		shutdown:          shutdown,
		processingChannel: make(chan bool, maxWorkers),
		finishedChannel:   make(chan worker.ID, maxWorkers),
		activeWorkers:     map[worker.ID]carbon.Carbon{},
		currentTransition: initialTransition,
	}
}

func (workerPool *WorkerPool) Start(actionHandler func(worker worker.Worker)) {
	for workerPool.currentTransition.state != stateExit {
		switch workerPool.currentTransition.state {
		case stateInitial:
			workerPool.processInitialState()
		case stateMain:
			workerPool.processMainState()
		case stateWait:
			workerPool.processWaitState()
		case stateSleep:
			workerPool.processSleepState()
		case stateLaunch:
			workerPool.processLaunchState(actionHandler)
		case stateQuit:
			workerPool.processQuitState()
		case stateTimeout:
			workerPool.processTimeoutState()
		case stateFinish:
			workerPool.processFinishState()
		case stateProcessing:
			workerPool.processProcessingState()
		default:
			panic("invalid state transition")
		}
	}
	workerPool.shutdown.doneChannel <- true
}

func (workerPool *WorkerPool) processLaunchState(actionHandler func(worker worker.Worker)) {
	newWorker := worker.NewWorker(workerPool.processingChannel, workerPool.finishedChannel)
	workerPool.activeWorkers[newWorker.ID] = *carbon.Now()

	go actionHandler(newWorker)

	workerPool.goToState(stateWait, nil)
}

func (workerPool *WorkerPool) processWaitState() {
	select {
	case finished := <-workerPool.finishedChannel:
		payload := &payload{
			workerID: finished,
		}
		workerPool.goToState(stateFinish, payload)
	case processing := <-workerPool.processingChannel:
		payload := &payload{
			isProcessing: processing,
		}
		workerPool.goToState(stateProcessing, payload)
	case <-workerPool.shutdown.initiateChannel:
		workerPool.goToState(stateQuit, nil)
	case <-time.After(time.Second * time.Duration(workerPool.timeout)):
		workerPool.goToState(stateTimeout, nil)
	}
}
