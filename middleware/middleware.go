/* 
Package middleware provides shared middleware that is used in HTTP services.
*/
package middleware

import (
	"fmt"
	"net/http"
)

// CORS handling middleware
type CorsHandler struct {
	corsHostAndPort string
	delegate        http.Handler
}

// Construct a new CORS handler.
//
// corsHostAndPort is a string representation of the allowed origin.
func NewCorsHandler(corsHostAndPort string, handler http.Handler) *CorsHandler {
	return &CorsHandler{
		corsHostAndPort: corsHostAndPort,
		delegate:        handler,
	}
}

/*
Write the appropriate CORS headers required for eventbus to function
properly. Currently two headers are added:

 * Access-Control-Allow-Origin: http://host:port
 * Access-Control-Allow-Headers: Content-Type

*/
func (handler *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://%s", handler.corsHostAndPort))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	handler.delegate.ServeHTTP(w, r)
}
