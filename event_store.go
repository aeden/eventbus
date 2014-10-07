package eventbus

type EventStore interface {
	WriteEvent(*Event) error
}

type InMemoryEventStore struct {
	Events []*Event
}

func NewInMemoryEventStore() EventStore {
	return &InMemoryEventStore{}
}

func (store *InMemoryEventStore) WriteEvent(event *Event) (err error) {
	store.Events = append(store.Events, event)
	return
}
