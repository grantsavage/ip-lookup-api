package dns

import (
	"net"
	"testing"
)

func TestLookupIP(t *testing.T) {
	type input struct {
		err       error
		ipAddress string
		response  []string
	}
	type want struct {
		err          error
		responseCode string
	}

	tests := []struct {
		description string
		input       input
		want        want
	}{
		{
			description: "should return result on successful lookup",
			input: input{
				ipAddress: "1.2.3.4",
				response:  []string{"127.0.0.4"},
			},
			want: want{
				responseCode: "127.0.0.4",
			},
		},
		{
			description: "should return error if lookup failed",
			input: input{
				ipAddress: "1.2.3.4",
				err:       net.ErrClosed,
			},
			want: want{
				err: net.ErrClosed,
			},
		},
		{
			description: "should return error if no response was returned",
			input: input{
				ipAddress: "1.2.3.4",
				response:  []string{},
			},
			want: want{
				err: ErrorNoResponse,
			},
		},
		{
			description: "should return error if response does not match expected format",
			input: input{
				ipAddress: "1.2.3.4",
				response:  []string{"123.4.5.6"},
			},
			want: want{
				err: ErrorUnexpectedResponse,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			lookupFunc := func(string) ([]string, error) {
				return test.input.response, test.input.err
			}
			got, err := LookupIP(net.ParseIP(test.input.ipAddress), "zen.spamhaus.org", lookupFunc)

			// Check error condition
			assertError(t, err, test.want.err)

			// Check response codes
			if test.want.responseCode != "" && got != nil && test.want.responseCode != got.String() {
				t.Errorf("got %s, want %s", got, test.want.responseCode)
			}
		})
	}
}

func TestSearchIPBlocklist(t *testing.T) {
	type input struct {
		err       error
		ipAddress string
		response  []string
	}
	type want struct {
		err          error
		responseCode string
	}

	tests := []struct {
		description string
		input       input
		want        want
	}{
		{
			description: "should return response code",
			input: input{
				ipAddress: "1.2.3.4",
				response:  []string{"127.0.0.4"},
			},
			want: want{
				responseCode: "127.0.0.4",
			},
		},
		{
			description: "should return error when lookup fails",
			input: input{
				err: net.ErrClosed,
			},
			want: want{
				err: net.ErrClosed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			lookupFunc := func(string) ([]string, error) {
				return test.input.response, test.input.err
			}
			got, err := SearchIPBlocklist(net.ParseIP(test.input.ipAddress), lookupFunc)

			// Check error condition
			assertError(t, err, test.want.err)

			// Check response codes
			if test.want.responseCode != "" && got != nil && test.want.responseCode != got.String() {
				t.Errorf("got %s, want %s", got, test.want.responseCode)
			}
		})
	}
}
