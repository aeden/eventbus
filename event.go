package eventbus

type Event struct {
	Name    string            `json:"name"`
	Data    interface{}       `json:"data"`
	Context map[string]string `json:"context"`
}

func NewEvent() *Event {
	return &Event{}
}
