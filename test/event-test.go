package main

import (
	"github.com/zlbbq/go-event"
	"github.com/zlbbq/go-logger"
	"time"
)

type Obj struct {
	EventA *event.Event
	Events *event.Events
}

type Arg struct {
	Name string
	Age int
}

func main() {
	obj := Obj{}
	obj.EventA = event.NewEvent("A")
	obj.Events = event.CreateEvents()

//	obj.EventA.AddListener(func(arg event.EventArgument){
//		if(arg != nil) {
//			logger.Info("name 1 -> %s", arg.(string))
//		}
//	})

	l1 := obj.EventA.AddListener(func(arg event.EventArgument){
		if(arg != nil) {
			logger.Info("name 3 -> %s", arg.(string))
		}
	})

	obj.Events.On("eventName1", func(arg event.EventArgument) {
		logger.Info("name 0 -> %s", arg.(Arg).Name)
		logger.Info("age 0 -> %d", arg.(Arg).Age)
	})

	obj.Events.On("eventName1", func(arg event.EventArgument) {
		logger.Info("name 2 -> %s", arg.(Arg).Name)
		logger.Info("age 2 -> %d", arg.(Arg).Age)
	})

	arg := Arg{
		"ZHANGLEI",
		1,
	}

	go (func(){
		for {
			time.Sleep(1 * time.Second)
			obj.EventA.Trigger("ZLBBQ")

			time.Sleep(1 * time.Second)
			obj.Events.Trigger("eventName1", arg)

			obj.EventA.RemoveListener(l1) // "name 3" print only once
		}
	})()

	go (func() {
		// print nothing 5 seconds later
		time.Sleep(5 * time.Second)
		obj.Events.Off("eventName1", nil)
	})()

	for {}

}
