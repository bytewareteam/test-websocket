package websocket

import "sync"

type (
	EventListenerData struct {
		Session *Session
		Event   *Event
	}
	EventChanManager map[string][] chan EventListenerData
	EventListener    struct {
		Events EventChanManager
		Mutex  *sync.Mutex
	}
)

func NewEventListener() *EventListener {
	return &EventListener{
		Events: make(EventChanManager),
		Mutex:  &sync.Mutex{},
	}
}

func (el *EventListener) Subscribe(event string, ch chan EventListenerData) {
	el.Mutex.Lock()
	if _, ok := el.Events[event]; ok {
		el.Events[event] = append(el.Events[event], ch)
	} else {
		el.Events[event] = []chan EventListenerData{ch}
	}
	el.Mutex.Unlock()
}

func (el *EventListener) Unsubscribe(event string, ch chan EventListenerData) {
	if _, ok := el.Events[event]; ok {
		for i := range el.Events[event] {
			if el.Events[event][i] == ch {
				el.Events[event] = append(el.Events[event][:i], el.Events[event][i+1:]...)
				break
			}
		}
	}
}

func (el *EventListener) Next(event *Event, session *Session) {
	if _, ok := el.Events[event.Name]; ok {
		for _, handler := range el.Events[event.Name] {
			go func(handler chan EventListenerData) {
				handler <- EventListenerData{session, event}
			}(handler)
		}
	}
}
