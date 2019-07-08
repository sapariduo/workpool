package workerpool

import (
	"testing"
	"time"
)

func TestShutdown_WaitForSignal(t *testing.T) {
	type fields struct {
		initiateChannel chan bool
		doneChannel     chan bool
		timeout         time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdown := &Shutdown{
				initiateChannel: tt.fields.initiateChannel,
				doneChannel:     tt.fields.doneChannel,
				timeout:         tt.fields.timeout,
			}
			shutdown.WaitForSignal()
		})
	}
}
