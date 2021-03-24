package services

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// ReverseIP reverses the given IP
func ReverseIP(ip net.IP) (net.IP, error) {
	// If ip is not an IPv4 address, return empty
	if ip.To4() == nil {
		return nil, errors.New("Provided IP " + ip.String() + " is not an IPv4 address.")
	}

	// Split address by address delimeter
	splitAddress := strings.Split(ip.String(), ".")

	// Reverse address parts
	for i, j := 0, len(splitAddress)-1; i < len(splitAddress)/2; i, j = i+1, j-1 {
		splitAddress[i], splitAddress[j] = splitAddress[j], splitAddress[i]
	}

	// Join the reversed address parts
	reversedAddress := strings.Join(splitAddress, ".")

	return net.ParseIP(reversedAddress), nil
}

// LookupIP looks up the target IP against the DNSBL address
func LookupIP(targetIP net.IP, dnsblAddress string) (net.IP, error) {
	// Create lookup address from target IP and DNSBL address
	lookup := fmt.Sprintf("%s.%s", targetIP.String(), dnsblAddress)

	// Perform lookup
	response, err := net.LookupHost(lookup)
	if err != nil {
		return nil, err
	}

	// Check the response length
	if len(response) == 0 {
		return nil, errors.New("no response from address lookup")
	}

	// We want just the first result from the response
	ip := response[0]
	match, err := regexp.MatchString("^127.0.0.*", ip)
	if err != nil {
		return nil, err
	}

	// Check if response is valid
	if !match {
		return nil, errors.New("response did not match expected response code")
	}

	return net.ParseIP(ip), nil
}
