package eventbus

type ServicesConfig struct {
	Services []map[string]string
}

type Authenticator interface {
	Authenticate(credentials interface{}) (bool, error)
}

type DefaultAuthenticator struct{}

func (authenticator *DefaultAuthenticator) Authenticate(credentials interface{}) (bool, error) {
	return true, nil
}
