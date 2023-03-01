package hermes

import "sync"

type Channel struct {
	name        string
	actors      []string
	mutexActors sync.Mutex
}

func (r *Channel) findActor(name string) (index int, found bool) {
	for i, actor := range r.actors {
		if actor == name {
			return i, true
		}
	}
	return -1, false
}

func (r *Channel) addActor(name string) {
	_, found := r.findActor(name)
	if !found {
		r.mutexActors.Lock()
		r.actors = append(r.actors, name)
		r.mutexActors.Unlock()
	}
}

func (r *Channel) removeActor(name string) {
	index, found := r.findActor(name)
	if found {
		r.mutexActors.Lock()
		r.actors = append(r.actors[:index], r.actors[index+1:]...)
		r.mutexActors.Unlock()
	}
}

func (r *Channel) broadcast(event Event) {
	index, found := r.findActor(event.actor)
	if found {
		event.setChannel(r.name)
		for i, actor := range r.actors {
			if i != index {
				go GetActor(actor).handle(event)
			}
		}
	}
}
