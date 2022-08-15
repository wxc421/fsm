package fsm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestFSM_Call(t *testing.T) {
	type fields struct {
		state           State
		name            string
		stateMachineMap map[State]map[EventName]*Action
		locker          sync.Mutex
	}
	type args struct {
		eventName EventName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    State
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:           tt.fields.state,
				name:            tt.fields.name,
				stateMachineMap: tt.fields.stateMachineMap,
				locker:          tt.fields.locker,
			}
			got, err := fsm.Call(tt.args.eventName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Call() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type Step1Observer struct{}

func (o Step1Observer) Receive(oldState, newState State, eventName EventName, data any, err error) {
	time.Sleep(time.Second * 5)
	fmt.Println("receive")
}

func TestFSM_RegisterStateMachine(t *testing.T) {
	const (
		step0 State = iota
		step1
		step2
		step3
	)
	const (
		event1 EventName = "event1"
		event2 EventName = "event2"
		event3 EventName = "event3"
		event4 EventName = "event4"
	)

	fsm := NewFSM("test", step0)
	fsm.RegisterStateMachine(step0, event1, func(data any) (State, error) {
		fmt.Println("step0 -> step1")
		return step1, nil
	})
	fsm.RegisterStateMachine(step1, event2, func(data any) (State, error) {
		fmt.Println("step1 -> step2")
		return step2, nil
	}, WithObservers(Step1Observer{}, Step1Observer{}, Step1Observer{}, Step1Observer{}), WithDrainWorker(8))
	fsm.RegisterStateMachine(step2, event3, func(data any) (State, error) {
		fmt.Println("step2 -> step3")
		return step3, nil
	})
	assert := assert.New(t)
	status, err := fsm.Call("event1")
	assert.Equal(status, step1)
	assert.Nil(err)

	status, err = fsm.Call("event2")
	assert.Equal(status, step2)
	assert.Nil(err)

	status, err = fsm.Call("event2")
	assert.Equal(status, step2)
	assert.NotNil(err)
}

func TestNewFSM(t *testing.T) {
	type args struct {
		name      string
		initState State
	}
	tests := []struct {
		name string
		args args
		want *FSM
	}{
		{"a", args{
			name:      "a",
			initState: 0,
		}, &FSM{
			state:           0,
			name:            "a",
			stateMachineMap: make(map[State]map[EventName]*Action),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFSM(tt.args.name, tt.args.initState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFSM() = %v, want %v", got, tt.want)
			}
		})
	}
}
