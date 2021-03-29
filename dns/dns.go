package dns

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/grantsavage/ip-lookup-api/db"
	"github.com/grantsavage/ip-lookup-api/graph/model"
	uuid "github.com/satori/go.uuid"
)

// HostLookupFunc is a function interface for performing the host lookup
type HostLookupFunc func(string) ([]string, error)

// Error definitions
var ErrorNoResponse = errors.New("no response from address lookup")
var ErrorUnexpectedResponse = errors.New("response did not match expected response code")

// LookupIP looks up the target IP against the DNSBL address
func LookupIP(targetIP net.IP, dnsblAddress string, lookupFunc HostLookupFunc) (net.IP, error) {
	// Create lookup address from target IP and DNSBL address
	lookup := fmt.Sprintf("%s.%s", targetIP.String(), dnsblAddress)

	// Perform lookup
	response, err := lookupFunc(lookup)
	if err != nil {
		return nil, err
	}

	// Check the response length
	if len(response) == 0 {
		return nil, ErrorNoResponse
	}

	// We want just the first result from the response
	ip := response[0]
	match, err := regexp.MatchString("^127.0.0.*", ip)
	if err != nil {
		return nil, err
	}

	// Check if response is valid
	if !match {
		return nil, ErrorUnexpectedResponse
	}

	return net.ParseIP(ip), nil
}

// SearchIPBlocklist normalizes the given IP and performs the blocklist lookup
func SearchIPBlocklist(ipAddress net.IP, lookupFunc HostLookupFunc) (net.IP, error) {
	// Reverse the IP
	reversedIp := ReverseIP(ipAddress)

	// Lookup the IP
	responseCode, err := LookupIP(reversedIp, "zen.spamhaus.org", lookupFunc)
	if err != nil {
		return nil, err
	}

	return responseCode, err
}

// BlocklistWorker loops over a list of IPs and additionally searches and stores the lookup results
func BlocklistWorker(ips []net.IP) {
	// Kick off a background task to lookup each valid IP
	for _, ipAddress := range ips {
		log.Printf("querying blocklist for IP address " + ipAddress.String())

		// Search IP blocklist and get response code
		responseCode, err := SearchIPBlocklist(ipAddress, net.LookupHost)
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

		// Upsert lookup result
		err = db.UpsertIPLookupResult(result)
		if err != nil {
			log.Printf("error occurred while storing result: %s\n", err.Error())
		}
	}
}
