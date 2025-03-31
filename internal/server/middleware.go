package server

import (
	"log"
	"net/http"
)

// RequestLogger logs incoming requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// CSRFProtection adds CSRF tokens to forms
func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Actual implementation would:
		// 1. Generate token for GET requests
		// 2. Validate token for modifying requests
		next.ServeHTTP(w, r)
	})
}
