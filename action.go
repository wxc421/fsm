package fsm

type (
	// Action is the core that wraps the basic Event methods.
	Action struct {
		eventName    EventName
		fn           EventFunc
		observers    []EventObserver
		drainWorkers int
		fsm          *FSM
	}

	// EventObserver is the interface.When the event is processed,
	//  it can notify the observers asynchronously and execute their own business.
	EventObserver interface {
		Receive(oldState, newState State, eventName EventName, param *Param, err error)
	}

	ActionOpt func(action *Action)
)

func (action *Action) run(param *Param) (State, error) {
	state, err := action.fn(param)
	runner := NewTaskRunner(action.drainWorkers)
	for _, observer := range action.observers {
		runner.Schedule(func() {
			observer.Receive(action.fsm.state, state, action.eventName, param, err)
		})
	}
	return state, err
}

// WithObservers adds observers to the event.
func WithObservers(observers ...EventObserver) ActionOpt {
	return func(action *Action) {
		if len(observers) == 0 {
			return
		}
		action.observers = append(action.observers, observers...)
	}
}

// WithDrainWorker limit goroutine number.
func WithDrainWorker(drainWorker int) ActionOpt {
	return func(action *Action) {
		if drainWorker <= 0 {
			drainWorker = 1
		}
		action.drainWorkers = drainWorker
	}
}
