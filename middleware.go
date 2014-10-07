package eventbus

import (
	"fmt"
	"log"
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

// Authorization middleware
type AuthorizationHandler struct {
	delegate http.Handler
}

func NewAuthorizationHandler(handler http.Handler) http.Handler {
	return &AuthorizationHandler{
		delegate: handler,
	}
}

func (handler *AuthorizationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header["Authorization"]
	log.Printf("Authorization: %s", authorization)
	handler.delegate.ServeHTTP(w, r)
}
