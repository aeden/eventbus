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
	servicesConfig *ServicesConfig
	delegate       http.Handler
}

func NewAuthorizationHandler(servicesConfig *ServicesConfig, handler http.Handler) http.Handler {
	return &AuthorizationHandler{
		servicesConfig: servicesConfig,
	}
}

func (handler *AuthorizationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	log.Printf("Authorization: %s", authorization)
	if authorization != "" {
		for _, serviceConfig := range handler.servicesConfig.Services {
			if serviceConfig["token"] == authorization {
				log.Printf("Authenticated service %s", serviceConfig["name"])
			}
		}
	}
	handler.delegate.ServeHTTP(w, r)
}
