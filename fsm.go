package fsm

import (
	"errors"
	"sync"
)

type (
	State     int32
	EventName string
	EventFunc func(data any) (State, error)
)

// FSM is the finite state machine
type FSM struct {
	// current state
	state           State
	name            string
	stateMachineMap map[State]map[EventName]*Action
	locker          sync.Mutex
}

// NewFSM creates a new fsm
func NewFSM(name string, initState State) *FSM {
	return &FSM{
		state:           initState,
		name:            name,
		stateMachineMap: make(map[State]map[EventName]*Action),
	}
}

// RegisterStateMachine register state machine
func (fsm *FSM) RegisterStateMachine(state State, eventName EventName, fn EventFunc, opts ...ActionOpt) {
	fsm.locker.Lock()
	defer fsm.locker.Unlock()

	if fsm.stateMachineMap[state] == nil {
		fsm.stateMachineMap[state] = make(map[EventName]*Action)
	}
	action := &Action{
		eventName:    eventName,
		fn:           fn,
		drainWorkers: 1,
		fsm:          fsm,
	}
	for _, opt := range opts {
		opt(action)
	}
	fsm.stateMachineMap[state][eventName] = action
}

// Call the state's event func
func (fsm *FSM) Call(eventName EventName, opts ...ParamOption) (State, error) {
	events, ok := fsm.stateMachineMap[fsm.state]
	if !ok || events == nil {
		return fsm.state, errors.New("can't find ")
	}

	action, ok := events[eventName]
	if !ok || action == nil {
		return fsm.state, errors.New("can't find ")
	}

	// call eventName func
	param := &Param{}
	for _, fn := range opts {
		fn(param)
	}
	state, err := action.run(param)
	if err != nil {
		return fsm.state, err
	}
	// update state
	fsm.state = state
	return fsm.state, nil
}

// AutoCall the state's event func multipart
func (fsm *FSM) AutoCall(eventNames ...EventName) (State, error) {
	for _, eventName := range eventNames {
		events, ok := fsm.stateMachineMap[fsm.state]
		if !ok || events == nil {
			return fsm.state, errors.New("can't find ")
		}

		action, ok := events[eventName]
		if !ok || action == nil {
			return fsm.state, errors.New("can't find ")
		}

		// call eventName func
		state, err := action.run(nil)
		if err != nil {
			return fsm.state, err
		}
		// update state
		fsm.state = state
	}
	return fsm.state, nil
}
