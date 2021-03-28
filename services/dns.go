package services

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/grantsavage/ip-lookup-api/db"
	"github.com/grantsavage/ip-lookup-api/graph/model"
	uuid "github.com/satori/go.uuid"
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

func SearchIPBlocklist(ipAddress net.IP) (net.IP, error) {
	// Reverse the IP
	reversedIp, err := ReverseIP(ipAddress)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reversing the IP %s: %s", ipAddress.String(), err.Error())
	}

	// Lookup the IP
	responseCode, err := LookupIP(reversedIp, "zen.spamhaus.org")
	if err != nil {
		return nil, fmt.Errorf("error occurred during IP lookup: %s", err.Error())
	}

	return responseCode, err
}

// ValidateIPs validates and normalizes a list of IPs
func ValidateIPs(ips []string) ([]net.IP, error) {
	validIPs := []net.IP{}

	// Validate and convert to native IP type
	for _, ipString := range ips {
		// Parse and validate the IP. If IP is not valid, return an error
		ip := net.ParseIP(ipString)
		if ip == nil {
			return nil, fmt.Errorf("provided IP %s is not a valid IP", ipString)
		}

		// If ip is not an IPv4 address, return an error
		if ip.To4() == nil {
			return nil, fmt.Errorf("provided IP %s is not an IPv4 address", ipString)
		}

		// If IP is valid, add it to list of IPs to lookup
		validIPs = append(validIPs, ip)
	}

	return validIPs, nil
}

// BlocklistWorker loops over a list of IPs and additionally searches and stores the lookup results
func BlocklistWorker(ips []net.IP) {
	// Kick off a background task to lookup each valid IP
	for _, ipAddress := range ips {
		log.Printf("querying blocklist for IP address " + ipAddress.String())

		// Search IP blocklist and get response code
		responseCode, err := SearchIPBlocklist(ipAddress)
		if err != nil {
			log.Printf("error occurred while searching IP blocklist: %s\n", err.Error())
			continue
		}

		log.Printf("storing result for IP " + ipAddress.String())

		// Bulid result
		result := model.IPLookupResult{
			UUID:         uuid.NewV4().String(),
			IPAddress:    ipAddress.String(),
			ResponseCode: responseCode.String(),
			CreatedAt:    time.Now().Format(time.RFC3339),
			UpdatedAt:    time.Now().Format(time.RFC3339),
		}

		// Store result
		err = db.StoreIPLookupResult(result)
		if err != nil {
			log.Printf("error occurred while storing result: %s\n", err.Error())
		}
	}
}
