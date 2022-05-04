package network

type Event int

const (
	LISTENING = iota
	NEWPEER
	MESSAGE
)

type Handler func(*Peer, ...any)
type Events map[Event]Handler

func (events Events) AddListener(e Event, handler Handler) {
	events[e] = handler
}

func (events Events) Emit(e Event, route *Peer, params ...any) {
	if event, ok := events[e]; ok {
		event(route, params...)
	}
}

func (events Events) RemoveListener(e Event) {
	delete(events, e)
}
