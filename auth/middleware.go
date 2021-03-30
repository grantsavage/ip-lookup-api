package auth

import (
	"encoding/json"
	"net/http"
	"os"
)

var configuredUsername string
var configuredPassword string

// init will read the authentication credentials from the environment
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
			// No need to check the error, as this is static
			response, _ := json.Marshal(map[string]string{
				"errors": "Unauthorized",
			})

			http.Error(w, string(response), http.StatusUnauthorized)
			return
		}

		// Continue if credentials are ok
		next.ServeHTTP(w, r)
	})
}
