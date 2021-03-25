package auth

import (
	"encoding/json"
	"net/http"
	"os"
)

var configuredUsername string
var configuredPassword string

func init() {
	configuredUsername = os.Getenv("AUTH_USERNAME")
	configuredPassword = os.Getenv("AUTH_PASSWORD")
}

// checkCredentials compares the supplied credentials to the application configured credentials
func checkCredentials(username, password string) bool {
	return username == configuredUsername && password == configuredPassword
}

// Middleware authenticates each request against the configured application credentials
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract username and password from Authorization header
		username, password, ok := r.BasicAuth()

		// If provided credentials do not match configured credentials, throw HTTP error
		if !ok || !checkCredentials(username, password) {
			response, err := json.Marshal(map[string]string{
				"errors": "Unauthorized",
			})
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			http.Error(w, string(response), http.StatusForbidden)
			return
		}

		// Continue if credentials are ok
		next.ServeHTTP(w, r)
	})
}
