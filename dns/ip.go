package dns

import (
	"errors"
	"net"
	"strings"
)

// Error definitions
var ErrorInvalidIP = errors.New("provided IP is not a valid IP")
var ErrorNonIPV4Address = errors.New("provided IP is not an IPv4 address")

// ReverseIP reverses the given IP
func ReverseIP(ip net.IP) net.IP {
	// Split address by address delimeter
	splitAddress := strings.Split(ip.String(), ".")

	// Reverse address parts
	for i, j := 0, len(splitAddress)-1; i < len(splitAddress)/2; i, j = i+1, j-1 {
		splitAddress[i], splitAddress[j] = splitAddress[j], splitAddress[i]
	}

	// Join the reversed address parts
	reversedAddress := strings.Join(splitAddress, ".")

	return net.ParseIP(reversedAddress)
}

// ValidateIPs validates and normalizes a list of IPs
func ValidateIPs(ips []string) ([]net.IP, error) {
	validIPs := []net.IP{}

	// Validate and convert to native IP type
	for _, ipString := range ips {
		// Parse and validate the IP. If IP is not valid, return an error
		ip := net.ParseIP(ipString)
		if ip == nil {
			return nil, ErrorInvalidIP
		}

		// If ip is not an IPv4 address, return an error
		if ip.To4() == nil {
			return nil, ErrorNonIPV4Address
		}

		// If IP is valid, add it to list of IPs to lookup
		validIPs = append(validIPs, ip)
	}

	return validIPs, nil
}
