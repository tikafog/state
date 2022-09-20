package state

import (
	"context"
)

type TriggerFunc func(ctx context.Context, trigger *Trigger)

type Trigger struct {
	// From is the state before the transition.
	From uint32
	// To is the state after the transition.
	To uint32
	// cancelFunc is called in case the event is canceled.
	cancelFunc func()
	// Err is an optional error that can be returned from a callback.
	Err error
}

func newTrigger(fn context.CancelFunc) *Trigger {
	return &Trigger{
		cancelFunc: fn,
	}
}

// Cancel can be called in before_<Trigger> or leave_<Trigger> to cancel the
// current transition before it happens. It takes an optional error, which will
// overwrite e.Err if set before.
func (t *Trigger) Cancel(err ...error) {
	if t.cancelFunc != nil {
		t.cancelFunc()
		t.cancelFunc = nil
	}

	if len(err) > 0 {
		t.Err = err[0]
	}
}

// IsCanceled can be called in before_<Trigger> or leave_<Trigger> to cancel the
// current transition before it happens. It takes an optional error, which will
// overwrite e.Err if set before.
// @receiver *Trigger
// @return bool
func (t *Trigger) IsCanceled() bool {
	return t.cancelFunc == nil
}

func (t *Trigger) key(name string) string {
	return triggerKey{0, name}.String()
}

func (t *Trigger) stateKey(name string) string {
	return triggerKey{t.From, name}.String()
}
