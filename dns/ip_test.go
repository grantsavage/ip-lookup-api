package dns

import (
	"net"
	"testing"
)

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil && want == nil {
		return
	}

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want.Error() {
		t.Errorf("got error %q, want %q", got, want)
	}
}

func TestReverseIP(t *testing.T) {
	type input struct {
		ipAddress string
	}
	type want struct {
		ipAddress string
	}

	tests := []struct {
		description string
		input       input
		want        want
	}{
		{
			description: "should reverse valid IP",
			input: input{
				ipAddress: "1.2.3.4",
			},
			want: want{
				ipAddress: "4.3.2.1",
			},
		},
		{
			description: "should reverse valid IP",
			input: input{
				ipAddress: "127.0.0.1",
			},
			want: want{
				ipAddress: "1.0.0.127",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := ReverseIP(net.ParseIP(test.input.ipAddress))

			// Check response codes
			if test.want.ipAddress != "" && got != nil && test.want.ipAddress != got.String() {
				t.Errorf("got %s, want %s", got, test.want.ipAddress)
			}
		})
	}
}

func TestValidateIPs(t *testing.T) {
	type input struct {
		ipAddresses []string
	}
	type want struct {
		err error
	}

	tests := []struct {
		description string
		input       input
		want        want
	}{
		{
			description: "should return error for invalid IP",
			input: input{
				ipAddresses: []string{"not an IP"},
			},
			want: want{
				err: ErrorInvalidIP,
			},
		},
		{
			description: "should return error for invalid IP",
			input: input{
				ipAddresses: []string{"127123.0123123.0.1"},
			},
			want: want{
				err: ErrorInvalidIP,
			},
		},
		{
			description: "should return list of valid IPs",
			input: input{
				ipAddresses: []string{"1.2.3.4", "127.0.0.1"},
			},
			want: want{},
		},
		{
			description: "should return error for IPv6 address",
			input: input{
				ipAddresses: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			},
			want: want{
				err: ErrorNonIPV4Address,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got, err := ValidateIPs(test.input.ipAddresses)

			// Check error
			assertError(t, err, test.want.err)

			// Check result
			if err == nil && len(got) != len(test.input.ipAddresses) {
				t.Errorf("got list of length %d, wanted %d", len(got), len(test.input.ipAddresses))
			}
		})
	}
}
