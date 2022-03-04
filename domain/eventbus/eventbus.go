package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

// EventBus returns default event bus instance
var EventBus = NewEventBus()

// EventListener is the signature of functions that can handle an Event.
type EventListener func(data interface{})

type eventListenerWrapper struct {
	async    bool
	Listener EventListener
}

// The eventBus allows publish-subscribe-style communication between components
// without requiring the components to explicitly register with one another (and thus be aware of each other)
// Inspired by Guava eventBus ; this is a more lightweight implementation.
type eventBus struct {
	mutex     *sync.RWMutex
	listeners map[string][]*eventListenerWrapper
}

// NewEventBus return a new eventBus
func NewEventBus() *eventBus {
	return &eventBus{new(sync.RWMutex), map[string][]*eventListenerWrapper{}}
}

// Subscribe adds an EventListener to be called when an event is posted.
func (e *eventBus) Subscribe(event interface{}, listener EventListener) {
	e.subscribe(event, listener, false)
}

// SubscribeAsync adds an EventListener to be async called when an event is posted.
func (e *eventBus) SubscribeAsync(event interface{}, listener EventListener) {
	e.subscribe(event, listener, true)
}

func (e *eventBus) subscribe(event interface{}, listener EventListener, async bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	name := e.getEventName(event)
	list, ok := e.listeners[name]
	if !ok {
		list = []*eventListenerWrapper{}
	}
	list = append(list, &eventListenerWrapper{
		async,
		listener,
	})
	e.listeners[name] = list
}

// Publish sends an event to all subscribed listeners.
// Parameter data is optional ; Post can only have one map parameter.
func (e *eventBus) Publish(event interface{}) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	name := e.getEventName(event)
	list, ok := e.listeners[name]
	if !ok {
		return
	}
	for _, each := range list[:] { // iterate over unmodifyable copy
		if each.async {
			go each.Listener(event)
		} else {
			each.Listener(event)
		}
	}
}

// getEventName get the topic name by event type
func (e *eventBus) getEventName(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("%s/%s", t.PkgPath(), t.Name())
}
