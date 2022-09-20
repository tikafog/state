package state

import (
	"context"
	"sync"

	hmap "github.com/cornelk/hashmap"
	"go.uber.org/atomic"
)

type FSM struct {
	cond     *sync.Cond
	current  *atomic.Uint32
	triggers *hmap.Map[string, TriggerFunc]
}

func NewFSM() *FSM {
	return &FSM{
		cond:     sync.NewCond(NewStateLock()),
		current:  atomic.NewUint32(0),
		triggers: hmap.New[string, TriggerFunc](),
	}
}

// Current returns the current state of the FSM.
// @receiver *FSM
// @return uint32
func (f *FSM) Current() uint32 {
	return f.current.Load()
}

// Is returns true if state is the current state.
// @receiver *FSM
// @param uint32
// @return bool
func (f *FSM) Is(state uint32) bool {
	return state == f.current.Load()
}

// SetState allows the user to move to the given state from current state.
// // The call does not trigger any callbacks, if defined.
// @receiver *FSM
// @param uint32
func (f *FSM) SetState(state uint32) {
	f.current.Store(state)
}

// Trigger
// @receiver *FSM
// @param context.Context
// @param uint32
func (f *FSM) Trigger(ctx context.Context, state uint32) {
	f.TriggerWithFunc(ctx, state, nil)
}

// TriggerWithFunc
// @receiver *FSM
// @param context.Context
// @param uint32
func (f *FSM) TriggerWithFunc(ctx context.Context, state uint32, do func(ctx context.Context) error) {
	ctx, cancel := context.WithCancel(ctx)
	trigger := newTrigger(cancel)
	trigger.From = f.Current()
	trigger.To = state
	f.trigger(ctx, TriggerBefore, trigger)
	f.trigger(ctx, TriggerEnter, trigger)
	if do != nil {
		if err := do(ctx); err != nil {
			trigger.Cancel(err)
			f.trigger(ctx, TriggerFailed, trigger)
			return
		}
	}
	if !f.transition(f.Current(), state) {
		f.trigger(ctx, TriggerFailed, trigger)
		return
	}
	f.trigger(ctx, TriggerLeave, trigger)
	f.trigger(ctx, TriggerAfter, trigger)
}

// Callback
// @receiver *FSM
// @param context.Context
// @param uint32
func (f *FSM) Callback(name string, triggerFunc TriggerFunc) {
	f.triggers.Insert(triggerKey{0, name}.String(), triggerFunc)
}

// CallbackWithState
// @receiver *FSM
// @param context.Context
// @param uint32
func (f *FSM) CallbackWithState(name string, state uint32, triggerFunc TriggerFunc) {
	f.triggers.Insert(triggerKey{state, name}.String(), triggerFunc)
}

func (f *FSM) trigger(ctx context.Context, name string, trigger *Trigger) {
	t, ok := f.triggers.Get(trigger.stateKey(name))
	if ok {
		t(ctx, trigger)
	}
	t, ok = f.triggers.Get(trigger.key(name))
	if ok {
		t(ctx, trigger)
	}
}

func (f *FSM) transition(prev, next uint32) bool {
	return f.current.CompareAndSwap(prev, next)
}
