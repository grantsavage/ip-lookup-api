package services

import (
	"errors"
	"net"
	"testing"
)

func TestLookupIP(t *testing.T) {
	t.Run("should return lookup result", func(t *testing.T) {
		responseCode := "127.0.0.4"
		lookupFunc := func(string) ([]string, error) {
			return []string{responseCode}, nil
		}
		result, err := LookupIP(net.ParseIP("1.2.3.4"), "zen.spamhaus.org", lookupFunc)
		if err != nil {
			t.Errorf("returned error %s", err.Error())
		}
		if result.String() != responseCode {
			t.Errorf("got %s, want %s", result.String(), responseCode)
		}
	})

	t.Run("should return error if lookup failed", func(t *testing.T) {
		lookupFunc := func(string) ([]string, error) {
			return nil, errors.New("not connected to internet")
		}
		_, err := LookupIP(net.ParseIP("1.2.3.4"), "zen.spamhaus.org", lookupFunc)
		if err == nil {
			t.Errorf("did not return error")
		}
	})

	t.Run("should return error if no response was returned", func(t *testing.T) {
		lookupFunc := func(string) ([]string, error) {
			return []string{}, nil
		}
		_, err := LookupIP(net.ParseIP("1.2.3.4"), "zen.spamhaus.org", lookupFunc)
		if err == nil {
			t.Errorf("did not return error")
		}
	})

	t.Run("should return error if response does not match expected format", func(t *testing.T) {
		lookupFunc := func(string) ([]string, error) {
			return []string{"123.4.5.6"}, nil
		}
		_, err := LookupIP(net.ParseIP("1.2.3.4"), "zen.spamhaus.org", lookupFunc)
		if err == nil {
			t.Errorf("did not return error")
		}
	})
}

func TestSearchIPBlocklist(t *testing.T) {
	t.Run("should return response code", func(t *testing.T) {
		ipAddress := net.ParseIP("1.2.3.4")
		responseCode := "127.0.0.4"
		lookupFunc := func(string) ([]string, error) {
			return []string{"127.0.0.4"}, nil
		}
		got, err := SearchIPBlocklist(ipAddress, lookupFunc)
		if err != nil {
			t.Error("returned error", err.Error())
		}
		if got.String() != responseCode {
			t.Errorf("got %s, want %s", got, responseCode)
		}
	})

	t.Run("should return error when IP reversal fails", func(t *testing.T) {
		ipAddress := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		lookupFunc := func(string) ([]string, error) {
			return []string{"127.0.0.4"}, nil
		}
		_, err := SearchIPBlocklist(ipAddress, lookupFunc)
		if err == nil {
			t.Error("no error returned")
		}
	})

	t.Run("should return error when lookup fails", func(t *testing.T) {
		ipAddress := net.ParseIP("1.2.3.4")
		lookupFunc := func(string) ([]string, error) {
			return nil, errors.New("failed to lookup")
		}
		_, err := SearchIPBlocklist(ipAddress, lookupFunc)
		if err == nil {
			t.Error("no error returned")
		}
	})
}
