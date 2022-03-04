package eventbus

import (
	"fmt"
	"testing"
	"time"
)

type TestEvent struct{
}
type Test2Event struct{
}
func TestEventBus(t *testing.T){
	e := NewEventBus()
	e.SubscribeAsync(TestEvent{}, func(data interface{}) {
		t.Log(fmt.Sprintf("event1 async data %#v ",data))
	})
	e.Subscribe(TestEvent{}, func(data interface{}) {
		t.Log(fmt.Sprintf("event1 data %#v ",data))
	})
	e.Subscribe(Test2Event{}, func(data interface{}) {
		t.Log(fmt.Sprintf("event2 data %#v ",data))
	})
	e.Publish(TestEvent{})
	e.Publish(Test2Event{})
	time.Sleep(time.Second )
}
