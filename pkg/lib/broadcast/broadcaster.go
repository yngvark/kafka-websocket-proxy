// Package broadcast knows how to broadcast messages to subscribers
package broadcast

// Broadcaster is used for sending (broadcasting) messages to a number of subscribers
type Broadcaster struct {
	subscribers []chan<- string
}

// AddSubscriber adds a Subscriber to its list of subscribers
func (n *Broadcaster) AddSubscriber(subscriber chan<- string) {
	n.subscribers = append(n.subscribers, subscriber)
}

// BroadCast sends a message to all Subscriber-s
func (n *Broadcaster) BroadCast(msg string) error {
	for _, subscriber := range n.subscribers {
		subscriber <- msg
	}

	return nil
}

// New returns a new Broadcaster
func New() *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan<- string, 0),
	}
}
