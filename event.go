package event
import (
	"sync"
)

// Event struct
type Event struct {
	// Event Name
	Name      string
	/***************************************************************************/
	/***************************************************************************/
	// Coroutine channel
	triggerChan chan EventArgument
	// Listeners
	listeners []*Listener
	// Lock
	locker    sync.Locker
	// Max listeners
	maxListeners int
	// Is closed
	closed bool
}

type EventArgument interface{}

// Function pointer of a listener
type FnListener func(EventArgument)

type Listener struct {
	// Trigger only once
	once bool;
	// Element of listener list
	fn FnListener
}

const defaultMaxListeners = 10

// Create an event
func NewEvent(name string) *Event{
	e := new(Event)
	e.Name = name
	e.maxListeners = defaultMaxListeners
	e.triggerChan = make(chan EventArgument)
	e.listeners = make([]*Listener, 0)
	e.locker = &(sync.Mutex{})
	e.startListen()
	return e
}

// Trigger an event
func (e *Event) Trigger(arg EventArgument) {
	// Wait until event dispatching finished completely
	e.triggerChan <- arg
}

// Add a listener to an event
//
// Return nil if max listener size reached
func (e *Event) AddListener(fn FnListener) *Listener {
	return e.addListener(fn, false)
}

// Add a listener to an event but will be triggered only once
func (e *Event) Once(fn FnListener) {
	e.addListener(fn, true)
}

// Remove a listener from event
func (e *Event) RemoveListener(listener *Listener) bool {
	if(listener == nil) {
		return false
	}
	e.locker.Lock()
	defer e.locker.Unlock()
	idx := -1;
	for i, l := range e.listeners {
		if(l == listener) {
			idx = i
			break
		}
	}

	if(idx >= 0) {
		e.listeners = append(e.listeners[:idx], e.listeners[idx+1:]...)
		return true
	}
	return false
}

// Remove all listeners
func (e *Event) RemoveAllListeners() {
	e.locker.Lock()
	defer e.locker.Unlock()
	e.listeners = make([]*Listener, 0)
}

// Get num of listeners
func (e *Event) GetListenerNum() int {
	e.locker.Lock()
	defer e.locker.Unlock()
	return len(e.listeners)
}

// Set max listeners, default max listener is 10
//
// Call this function before any listener is add to event, or nothing happened
func (e *Event) SetMaxListeners(n int) {
	if(len(e.listeners) > 0) {
		return
	}

	if(n < 1) {
		n = defaultMaxListeners
	}
	e.maxListeners = n
}

/***************************************************************************/
// Start a goroutine to listen to trigger channel
func (e *Event) startListen() {
	go (func() {
		for {
			arg, ok := <- e.triggerChan
			if (ok == false) {
				// Channel closed
				break
			}
			e.dispatch(arg)
		}
	})()
}

// Dispatch event
func (e *Event)
dispatch(arg EventArgument) {
	for i:=0;i<len(e.listeners);i++ {
		l := e.listeners[i]
		go (func(listener *Listener) {
			// Call event handle function
			listener.fn(arg)
			if(listener.once == true) {
				go (func(){
					e.RemoveListener(listener)
				})()
			}
		})(l)
	}
}

// Private add listener
func (e *Event) addListener(fn FnListener, once bool) *Listener{
	e.locker.Lock()
	defer e.locker.Unlock()
	if(len(e.listeners) >= e.maxListeners) {
		return nil
	}
	l := &(Listener{
		once : once,
		fn : fn,
	})
	e.listeners = append(e.listeners, l)
	return l
}

