# 状态机

## Usage:
```go
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
```

## link
- https://www.cnblogs.com/prelude1214/p/16055915.html