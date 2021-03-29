package auth

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckCredentials(t *testing.T) {
	configuredUsername = "test"
	configuredPassword = "test"

	type input struct {
		username string
		password string
	}

	tests := []struct {
		description string
		input       input
		want        bool
	}{
		{
			description: "should return true if credentials match",
			input: input{
				username: "test",
				password: "test",
			},
			want: true,
		},
		{
			description: "should return false if credentials do not match",
			input: input{
				username: "test",
				password: "wrong",
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run("should return true if credentials match", func(t *testing.T) {
			got := checkCredentials(test.input.username, test.input.password)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	configuredUsername = "test"
	configuredPassword = "test"

	type input struct {
		username string
		password string
	}

	tests := []struct {
		description string
		input       input
		want        int
	}{
		{
			description: "should return ok on authenticated request",
			input: input{
				username: "test",
				password: "test",
			},
			want: http.StatusOK,
		},
		{
			description: "shouuld return unauthorized status code on unauthenticated request",
			input: input{
				username: "test",
				password: "wrongpassword",
			},
			want: http.StatusUnauthorized,
		},
	}

	// Create handler to send fake successful response
	nextHandler := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("ok"))
	})

	// Attach middleware to test handler
	handler := Middleware(nextHandler)

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// Create request
			request, _ := http.NewRequest(http.MethodGet, "http://testing", nil)

			// Compute authorization header
			authToken := b64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", test.input.username, test.input.password)))
			authHeader := fmt.Sprintf("Basic %s", authToken)
			request.Header.Add("Authorization", authHeader)

			// Create response recorder to capture response from the request
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			// Check response status code
			statusCode := responseRecorder.Result().StatusCode
			if statusCode != test.want {
				t.Errorf("got status %d, want %d", statusCode, test.want)
			}
		})
	}
}
