package worker

import (
	"reflect"
	"testing"
)

func TestNewWorker(t *testing.T) {
	type args struct {
		processingChannel chan<- bool
		finishedChannel   chan<- ID
	}
	tests := []struct {
		name string
		args args
		want Worker
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWorker(tt.args.processingChannel, tt.args.finishedChannel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}
