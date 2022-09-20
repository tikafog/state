package state

import (
	"fmt"
)

var (
	TriggerBefore = "before"
	TriggerEnter  = "enter"
	TriggerLeave  = "leave"
	TriggerAfter  = "after"
	TriggerFailed = "failed"
)

type triggerKey struct {
	current uint32
	trigger string
}

func (t triggerKey) String() string {
	return fmt.Sprintf("%v(%d)", t.trigger, t.current)
}

//
//func TriggerLeave(current uint32) triggerKey {
//	return triggerKey{
//		current: current,
//		trigger: triggerLeave,
//	}
//}
//
//func TriggerEnter(current uint32) triggerKey {
//	return triggerKey{
//		current: current,
//		trigger: triggerEnter,
//	}
//}
//
//func TriggerAfter(current uint32) triggerKey {
//	return triggerKey{
//		current: current,
//		trigger: triggerAfter,
//	}
//}
//
//func TriggerBefore(current uint32) triggerKey {
//	return triggerKey{
//		current: current,
//		trigger: triggerBefore,
//	}
//}
//
//func TriggerFailed(current uint32) triggerKey {
//	return triggerKey{
//		current: current,
//		trigger: triggerFailed,
//	}
//}
