package eventbus

/*
Authenticators are responsible for determining if a service is allowed
to connect to the event bus and receive all incoming events.
*/
type Authenticator interface {
	Authenticate(credentials interface{}) (bool, error)
}

// A simple authenticator that allows any service to connect.
type AllowAllAuthenticator struct{}

// Always return true, thus allowing any service to connect.
func (authenticator *AllowAllAuthenticator) Authenticate(credentials interface{}) (bool, error) {
	return true, nil
}
