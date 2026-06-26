// Package events provides a small in-process pub/sub broker used to fan out
// real-time events (check completed, alert fired) to SSE subscribers.
package events

import "sync"

// Event is a single real-time event delivered to subscribers.
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// Publisher is the write side of the broker. Producers (scheduler, alert
// manager) depend on this narrow interface so they don't import the broker.
type Publisher interface {
	Publish(event Event)
}

// subBuffer is the per-subscriber channel buffer. Slow consumers that fill it
// have events dropped rather than blocking the publisher.
const subBuffer = 32

// Broker is an in-memory fan-out hub. Safe for concurrent use.
type Broker struct {
	mu   sync.RWMutex
	subs map[chan Event]struct{}
}

// NewBroker creates an empty broker.
func NewBroker() *Broker {
	return &Broker{subs: make(map[chan Event]struct{})}
}

// Subscribe registers a new subscriber and returns its event channel together
// with an unsubscribe function. The unsubscribe function is idempotent and must
// be called (e.g. via defer) when the subscriber goes away; it closes the
// channel after removing it from the subscriber set.
func (b *Broker) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, subBuffer)

	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs, ch)
			b.mu.Unlock()
			close(ch)
		})
	}

	return ch, unsubscribe
}

// Publish delivers the event to all current subscribers without blocking. A
// subscriber whose buffer is full drops the event. Publish takes the read lock
// and unsubscribe takes the write lock, so a channel is never sent on after it
// has been closed.
func (b *Broker) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subs {
		select {
		case ch <- event:
		default:
		}
	}
}

// SubscriberCount returns the number of active subscribers (for diagnostics).
func (b *Broker) SubscriberCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs)
}
