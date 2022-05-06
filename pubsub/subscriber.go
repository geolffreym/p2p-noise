package pubsub

import "sync"

type Subscriber struct {
	message chan *Message
	topics  map[Event]bool
	// closed  chan bool
	mutex sync.Mutex
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		message: make(chan *Message),
		topics:  map[Event]bool{},
	}
}

func (s *Subscriber) Emit(msg *Message) {
	// Get lock to enforce sync order messages
	// https://gobyexample.com/mutexes
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.message <- msg
}

func (s *Subscriber) Listen(cb Observer) {
	go func(call Observer) {
		for {
			if msg, ok := <-s.message; ok {
				keepAlive := call(msg)
				if keepAlive == false {
					return
				}
			}
		}
	}(cb)
}
