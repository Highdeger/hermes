package hermes

import (
	"sync"
)

type Actor struct {
	name          string
	input         chan Event
	channels      []string
	handlers      map[string]func(event Event)
	mutexChannels sync.Mutex
	mutexHandlers sync.Mutex
}

func (r *Actor) Handle(name string, handler func(event Event)) {
	r.mutexHandlers.Lock()
	r.handlers[name] = handler
	r.mutexHandlers.Unlock()
}

func (r *Actor) RemoveHandler(name string) {
	r.mutexHandlers.Lock()
	delete(r.handlers, name)
	r.mutexHandlers.Unlock()
}

func (r *Actor) findChannel(name string) (index int, found bool) {
	for i, channel := range r.channels {
		if channel == name {
			return i, true
		}
	}
	return -1, false
}

func (r *Actor) handle(event Event) {
	handler, found := r.handlers[event.name]
	if found {
		handler(event)
	}
}

func (r *Actor) Subscribe(channel string) {
	_, found := r.findChannel(channel)
	if !found {
		GetChannel(channel).addActor(r.name)
		r.mutexChannels.Lock()
		r.channels = append(r.channels, channel)
		r.mutexChannels.Unlock()
	}
}

func (r *Actor) Unsubscribe(channel string) {
	index, found := r.findChannel(channel)
	if found {
		GetChannel(channel).removeActor(r.name)
		r.mutexChannels.Lock()
		r.channels = append(r.channels[:index], r.channels[index+1:]...)
		r.mutexChannels.Unlock()
	}
}

func (r *Actor) Broadcast(channel string, event Event) {
	_, found := r.findChannel(channel)
	if found {
		event.setActor(r.name)
		GetChannel(channel).broadcast(event)
	}
}
