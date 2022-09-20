package state

import (
	"context"
	"fmt"
	"testing"
)

func TestNewFSM(t *testing.T) {
	tests := []struct {
		name string
		want *FSM
	}{
		// TODO: Add test cases.
		{
			name: "",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFSM()

			got.CallbackWithState(TriggerAfter, 7, func(ctx context.Context, name string, trigger *Trigger) {
				fmt.Println("triggered", name)
			})
			got.Trigger(context.TODO(), 7)
		})
	}
}
