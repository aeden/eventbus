package eventbus

/*
Authenticators are responsible for determining if a service is allowed
to connect to the event bus and receive all incoming events.
*/
type authenticator interface {
	authenticate(credentials interface{}) (bool, error)
}

// A simple authenticator that allows any service to connect.
type allowAllAuthenticator struct{}

// Always return true, thus allowing any service to connect.
func (authenticator *allowAllAuthenticator) authenticate(credentials interface{}) (bool, error) {
	return true, nil
}
