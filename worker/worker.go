package worker

import (
	uuid "github.com/satori/go.uuid"
)

type ID struct {
	uuid.UUID
}
type Worker struct {
	ID                ID
	ProcessingChannel chan<- bool
	FinishedChannel   chan<- ID
}

func NewWorker(processingChannel chan<- bool, finishedChannel chan<- ID) Worker {
	return Worker{
		ID:                generateID(),
		ProcessingChannel: processingChannel,
		FinishedChannel:   finishedChannel,
	}
}
func generateID() ID {
	return ID{UUID: uuid.NewV4()}
}
