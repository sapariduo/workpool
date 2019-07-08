package consumer

import "github.com/sapariduo/workpool/worker"

func (queue Queue) Consume(worker worker.Worker) error {
	// signal the pool that the worker has finished after the method is executed
	defer func() { worker.FinishedChannel <- worker.ID }()

	// poll the queue
	message, err := queue.service.Read(queue.name)
	if err != nil {
		// signal the pool there's no message being processed
		worker.ProcessingChannel <- false

		return err
	}

	if len(message.Messages) == 0 {
		// signal the pool there's no message being processed
		worker.ProcessingChannel <- false

		return nil
	}
	// signal the pool that a message will be processed
	worker.ProcessingChannel <- true

	// handle the message...
}
