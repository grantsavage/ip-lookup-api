package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/grantsavage/ip-lookup-api/db"
	"github.com/grantsavage/ip-lookup-api/graph/generated"
	"github.com/grantsavage/ip-lookup-api/graph/model"
	"github.com/grantsavage/ip-lookup-api/services"
	uuid "github.com/satori/go.uuid"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ips []string) ([]string, error) {
	log.Printf("Mutation.Enqueue invoked for %d IP(s)", len(ips))

	// Check length of IPs input
	if len(ips) > 100 {
		return nil, fmt.Errorf("provided list of IPs is too large. Max number of IPs is %d", len(ips))
	}

	// Stores slice of valid IPs
	validIps := []net.IP{}

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
		validIps = append(validIps, ip)
	}

	// Kick off a background task to lookup each valid IP
	for _, validIp := range validIps {
		go func(ipAddress net.IP) {
			// Reverse the IP
			reversedIp, err := services.ReverseIP(ipAddress)
			if err != nil {
				log.Printf("error occurred while reversing the IP %s: %s", ipAddress.String(), err.Error())
				return
			}

			// Lookup the IP
			responseCode, err := services.LookupIP(reversedIp, "zen.spamhaus.org")
			if err != nil {
				log.Printf("error occurred during IP lookup: %s", err.Error())
				return
			}

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
		}(validIp)
	}

	return ips, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPLookupResult, error) {
	log.Printf("Query.GetIPDetails invoked for IP: %s", ip)

	// Validate IP input
	validIp := net.ParseIP(ip)
	if validIp == nil {
		return nil, errors.New("Provided IP " + ip + " is not a valid IP.")
	}

	// Retrieve the lookup result from the database
	result, err := db.GetIPLookupResult(validIp)
	if err != nil {
		return nil, errors.New("error occurred while retrieving the IP lookup result: " + err.Error())
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
