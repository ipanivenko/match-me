package graph

import (
	"sync"

	"matchme-server/graph/model"
)

// PubSub is our in-memory "hub" for chat messages.
type PubSub struct {
	mu   sync.Mutex
	subs map[string][]chan *model.Message
}

// NewPubSub creates a new PubSub object.
func NewPubSub() *PubSub {
	return &PubSub{
		subs: make(map[string][]chan *model.Message),
	}
}

// Subscribe adds a new listener for a given topic (e.g., "chat:123").
// It returns a channel to receive messages on, and a function to call to unsubscribe.
func (ps *PubSub) Subscribe(topic string) (<-chan *model.Message, func()) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Create a new channel for this subscriber
	ch := make(chan *model.Message, 1)

	// Add the channel to the list of subscribers for this topic
	ps.subs[topic] = append(ps.subs[topic], ch)

	// Create a cleanup function to remove the channel
	unsubscribe := func() {
		ps.mu.Lock()
		defer ps.mu.Unlock()

		// Find the channel in the slice
		subs := ps.subs[topic]
		for i, sub := range subs {
			if sub == ch {
				// Remove it by slicing
				ps.subs[topic] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		close(ch)
	}
	
	return ch, unsubscribe
}

// Publish sends a message to all subscribers of a given topic.
func (ps *PubSub) Publish(topic string, msg *model.Message) {
	ps.mu.Lock()
	
	// Get all subscribers for this topic
	listeners := ps.subs[topic]
	
	// Unlock *before* sending, so we don't block
	ps.mu.Unlock() 

	// Send the message to each listener
	for _, ch := range listeners {
		// Use a non-blocking send in case a listener is slow
		select {
		case ch <- msg:
		default:
		
		}
	}
}