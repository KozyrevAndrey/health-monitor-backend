package events

import (
	"testing"
)

func TestBroker_PublishToSubscriber(t *testing.T) {
	b := NewBroker()
	ch, unsub := b.Subscribe()
	defer unsub()

	if b.SubscriberCount() != 1 {
		t.Fatalf("Expected 1 subscriber, got %d", b.SubscriberCount())
	}

	b.Publish(Event{Type: "check", Data: map[string]string{"id": "x"}})

	ev := <-ch
	if ev.Type != "check" {
		t.Errorf("Expected event type 'check', got %q", ev.Type)
	}
}

func TestBroker_FanOut(t *testing.T) {
	b := NewBroker()
	ch1, unsub1 := b.Subscribe()
	defer unsub1()
	ch2, unsub2 := b.Subscribe()
	defer unsub2()

	b.Publish(Event{Type: "alert"})

	for i, ch := range []<-chan Event{ch1, ch2} {
		if ev := <-ch; ev.Type != "alert" {
			t.Errorf("subscriber %d: expected 'alert', got %q", i, ev.Type)
		}
	}
}

func TestBroker_Unsubscribe(t *testing.T) {
	b := NewBroker()
	ch, unsub := b.Subscribe()

	unsub()
	if b.SubscriberCount() != 0 {
		t.Errorf("Expected 0 subscribers after unsubscribe, got %d", b.SubscriberCount())
	}

	// Channel should be closed.
	if _, ok := <-ch; ok {
		t.Error("Expected channel to be closed after unsubscribe")
	}

	// Idempotent + publishing after unsubscribe must not panic.
	unsub()
	b.Publish(Event{Type: "check"})
}

func TestBroker_DropsWhenFull(t *testing.T) {
	b := NewBroker()
	_, unsub := b.Subscribe() // never drained
	defer unsub()

	// Publishing far more than the buffer must not block or panic.
	for i := 0; i < subBuffer*4; i++ {
		b.Publish(Event{Type: "check"})
	}
}
