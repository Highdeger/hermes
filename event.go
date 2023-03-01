package hermes

import "fmt"

type Event struct {
	channel string
	actor   string
	name    string
	params  map[string]interface{}
}

func NewEvent(name string, params map[string]interface{}) *Event {
	return &Event{
		channel: "",
		actor:   "",
		name:    name,
		params:  params,
	}
}

func (r *Event) setChannel(channel string) {
	r.channel = channel
}

func (r *Event) setActor(actor string) {
	r.actor = actor
}

func (r *Event) getParam(key string) interface{} {
	param, found := r.params[key]
	if !found {
		panic(fmt.Sprintf("param '%s' is missing. (event: %+v)", key, r))
	}
	return param
}

func (r *Event) GetString(key string) string {
	param, ok := r.getParam(key).(string)
	if !ok {
		panic(fmt.Sprintf("param '%s' is not string. (event: %+v)", key, r))
	}
	return param
}

func (r *Event) GetBool(key string) bool {
	param, ok := r.getParam(key).(bool)
	if !ok {
		panic(fmt.Sprintf("param '%s' is not bool. (event: %+v)", key, r))
	}
	return param
}

func (r *Event) GetInt(key string) int {
	param, ok := r.getParam(key).(int)
	if !ok {
		panic(fmt.Sprintf("param '%s' is not int. (event: %+v)", key, r))
	}
	return param
}

func (r *Event) GetFloat64(key string) float64 {
	param, ok := r.getParam(key).(float64)
	if !ok {
		panic(fmt.Sprintf("param '%s' is not float64. (event: %+v)", key, r))
	}
	return param
}
