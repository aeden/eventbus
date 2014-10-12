package eventbus

// Represents a single event fired by either a client or a service.
type Event struct {
	Name    string            `json:"name"`
	Data    interface{}       `json:"data"`
	Context map[string]string `json:"context"`
}

// Construct a new event and return a pointer to it.
func NewEvent() *Event {
	return &Event{}
}
