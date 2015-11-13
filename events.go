package event
import "sync"

type Events struct {
	locker *sync.Mutex
	eventMap map[string]*Event
}

// Create new Events instance
//
// An Events instance is implemented by map[string]*Event
func CreateEvents() *Events{
	events := &(Events{})
	events.locker = &(sync.Mutex{})
	events.eventMap = make(map[string]*Event)
	return events
}

// Add a listener to an event
func (events *Events) On(name string, fn FnListener) *Listener {
	event := events.ensureEvent(name)
	return event.AddListener(fn)
}

// Remove a listener from an event
func (events *Events) Off(name string, listener *Listener) {
	event := events.ensureEvent(name)
	event.RemoveListener(listener)
}

// Add a listener to an event, refer Event.Once()
func (events *Events) Once(name string, fn FnListener) {
	event := events.ensureEvent(name)
	event.Once(fn)
}

// Trigger an event
func (events *Events) Trigger(name string, arg EventArgument) {
	event := events.ensureEvent(name)
	event.Trigger(arg)
}

/********************************************************************************************/
// Private ensure event instance is created
func (events *Events) ensureEvent(name string) *Event{
	events.locker.Lock()
	defer events.locker.Unlock()
	e := events.eventMap[name]
	if(e == nil) {
		e = NewEvent(name)
		events.eventMap[name] = e
	}
	return e
}