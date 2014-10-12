package eventbus

// Event stores are used to persist events when they are created.
type EventStore interface {
	// Write the event to the event store. Returns an error
	// if the event fails to write.
	WriteEvent(*Event) error
}

// Event store that simply discards events.
type NullEventStore struct {
}

func NewNullEventStore() *NullEventStore {
	return &NullEventStore{}
}

func (store *NullEventStore) WriteEvent(event *Event) (err error) {
	return
}

// A simple in-memory event store that puts all events in an Array.
type InMemoryEventStore struct {
	Events []*Event
}

// Construct a new in-memory event store. Returns a pointer to the event store.
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{}
}

func (store *InMemoryEventStore) WriteEvent(event *Event) (err error) {
	store.Events = append(store.Events, event)
	return
}
