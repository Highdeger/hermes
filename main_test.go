package hermes

import (
	"fmt"
	"sync/atomic"
	"testing"
)

var (
	actorsCount = 100000
	eventsCount = 20000
)

func TestStress(t *testing.T) {
	actorList := make([]*Actor, actorsCount)
	for i := 0; i < actorsCount; i++ {
		actor := GetActor(fmt.Sprintf("actor-%d", i))
		if i%2 == 0 {
			actor.Subscribe("chan-1")
			actor.Subscribe("chan-2")
			actor.Subscribe("chan-3")
			actor.Subscribe("chan-4")
			actor.Subscribe("chan-5")
		} else {
			actor.Subscribe("chan-1")
			actor.Subscribe("chan-3")
			actor.Subscribe("chan-5")
		}
		actorList[i] = actor
	}

	t.Logf("%d actors have been created. even actors subscribed to 5 channels and odd ones subscribed to 3 channels.", actorsCount)

	eventsList := make([]Event, eventsCount)
	for i := 0; i < eventsCount; i++ {
		event := NewEvent(fmt.Sprintf("event-%d", i), map[string]interface{}{})
		eventsList[i] = *event
	}

	t.Logf("%d events have been created.", eventsCount)
	t.Log("adding handlers...")

	resultList := make(map[string]*int64)
	for _, actor := range actorList {
		for _, event := range eventsList {
			var a int64 = 0
			resultKey := fmt.Sprintf("%s_%s", actor.name, event.name)
			resultList[resultKey] = &a
			actor.Handle(event.name, func(event Event) {
				atomic.AddInt64(resultList[resultKey], 1)
			})
		}
	}

	t.Log("handlers have been added for all events to all the actors.")
	t.Log("broadcasting events...")

	for _, event := range eventsList {
		actorList[0].Broadcast("chan-1", event)
		actorList[0].Broadcast("chan-2", event)
		actorList[0].Broadcast("chan-3", event)
		actorList[0].Broadcast("chan-4", event)
		actorList[0].Broadcast("chan-5", event)
	}

	t.Logf("all %d events have been broadcasted to all 5 channels by the first actor who is subscribed to all of 5 channels.", eventsCount)
	t.Log("checking if all events have been received by related actors...")

	for i, actor := range actorList {
		for _, event := range eventsList {
			resultKey := fmt.Sprintf("%s_%s", actor.name, event.name)
			if i != 0 {
				if i%2 == 0 {
					if *resultList[resultKey] != 5 {
						t.Errorf("\tgot: %d, expected: %d\n", *resultList[resultKey], 5)
					}
				} else {
					if *resultList[resultKey] != 3 {
						t.Errorf("\tgot: %d, expected: %d\n", *resultList[resultKey], 3)
					}
				}
			}
		}
	}

	if !t.Failed() {
		t.Logf("the test was successfully passed and all %d events has been captured by all %d actors and their handlers.", eventsCount, actorsCount)
	}
}
