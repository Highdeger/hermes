package hermes

import "sync"

var (
	channels      map[string]*Channel
	actors        map[string]*Actor
	mutexChannels sync.Mutex
	mutexActors   sync.Mutex
)

func init() {
	channels = make(map[string]*Channel)
	actors = make(map[string]*Actor)
	mutexChannels = sync.Mutex{}
	mutexActors = sync.Mutex{}
}

func GetChannel(name string) *Channel {
	channel, found := channels[name]
	if !found {
		channel = &Channel{
			name:        name,
			actors:      make([]string, 0),
			mutexActors: sync.Mutex{},
		}
		mutexChannels.Lock()
		channels[name] = channel
		mutexChannels.Unlock()
	}
	return channel
}

func GetActor(name string) *Actor {
	actor, found := actors[name]
	if !found {
		actor = &Actor{
			name:          name,
			input:         make(chan Event),
			channels:      make([]string, 0),
			handlers:      make(map[string]func(event Event)),
			mutexChannels: sync.Mutex{},
			mutexHandlers: sync.Mutex{},
		}
		mutexActors.Lock()
		actors[name] = actor
		mutexActors.Unlock()
	}
	return actor
}
