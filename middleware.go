package eventbus

import (
	"fmt"
	"net/http"
)

// CORS handling middleware
type CorsHandler struct {
	corsHostAndPort string
	delegate        http.Handler
}

func NewCorsHandler(corsHostAndPort string, handler http.Handler) http.Handler {
	return &CorsHandler{
		corsHostAndPort: corsHostAndPort,
		delegate:        handler,
	}
}

func (handler *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://%s", handler.corsHostAndPort))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	handler.delegate.ServeHTTP(w, r)
}

// Authentication middleware
type AuthenticationHandler struct {
	delegate http.Handler
}

func NewAuthenticationHandler(handler http.Handler) http.Handler {
	return &AuthenticationHandler{
		delegate: handler,
	}
}

func (handler *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.delegate.ServeHTTP(w, r)
}
