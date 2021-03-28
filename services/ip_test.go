package services

import (
	"net"
	"testing"
)

func TestReverseIP(t *testing.T) {
	t.Run("should return error if address is IPv6", func(t *testing.T) {
		_, err := ReverseIP(net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
		if err == nil {
			t.Error("expected to throw error")
		}
	})

	t.Run("should reverse valid IPs", func(t *testing.T) {
		tests := []struct {
			input string
			want  string
		}{
			{input: "1.2.3.4", want: "4.3.2.1"},
			{input: "127.0.0.1", want: "1.0.0.127"},
		}

		for _, test := range tests {
			got, err := ReverseIP(net.ParseIP(test.input))
			if err != nil {
				t.Error("expected to not throw error")
			}
			if got.String() != test.want {
				t.Errorf("got %q want %q", got, test.want)
			}
		}
	})
}

func TestValidateIPs(t *testing.T) {
	t.Run("should return error for invalid IP", func(t *testing.T) {
		tests := []struct {
			input     string
			wantError bool
		}{
			{input: "not an IP", wantError: true},
			{input: "127123.0123123.0.1", wantError: true},
		}

		for _, test := range tests {
			_, err := ValidateIPs([]string{test.input})
			if err == nil && test.wantError {
				t.Error("expected to return error")
			}
		}
	})

	t.Run("should return valid IP list", func(t *testing.T) {
		ips := []string{"1.2.3.4", "127.0.0.1"}

		validIPs, err := ValidateIPs(ips)
		if err != nil {
			t.Error("expected to not return error")
		}
		if len(validIPs) != len(ips) {
			t.Error("expected to return list of IPs")
		}
	})
}
