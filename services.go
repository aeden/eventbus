package eventbus

// Configuration for a service that may attach.
type ServiceConfig struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}
