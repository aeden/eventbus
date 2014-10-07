package eventbus

import (
  "net/http"
  "fmt"
)

// CORS handling middleware
type CorsHandler struct {
	corsHostAndPort string
	delegate        http.Handler
}

func CorsServer(corsHostAndPort string, handler http.Handler) http.Handler {
	return &CorsHandler{
		corsHostAndPort: corsHostAndPort,
		delegate:        handler,
	}
}

func (server *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://%s", server.corsHostAndPort))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	server.delegate.ServeHTTP(w, r)
}
