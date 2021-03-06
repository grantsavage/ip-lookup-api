package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/grantsavage/ip-lookup-api/db"
	"github.com/grantsavage/ip-lookup-api/dns"
	"github.com/grantsavage/ip-lookup-api/graph/generated"
	"github.com/grantsavage/ip-lookup-api/graph/model"
)

// Enqueue looks up and stores the response codes of a given list of IPs
func (r *mutationResolver) Enqueue(ctx context.Context, ips []string) ([]string, error) {
	log.Printf("Mutation.Enqueue invoked for %d IP(s)", len(ips))

	// Validate the list of IPs
	validIPs, err := dns.ValidateIPs(ips)
	if err != nil {
		log.Printf("error while validating IP addresses: %s", err)
		return nil, err
	}

	/* Kick off a background worker to process IPs. Ideally if this system were
	 * to receive a high amount of requests, I would batch the IPs and kick off a worker
	 * for each batch of IPs to improve the concurrency. */
	go dns.BlocklistWorker(r.Database, validIPs)

	return ips, nil
}

// GetIPDetails fetches the lookup details of a given IP
func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPLookupResult, error) {
	log.Printf("Query.GetIPDetails invoked for IP: %s", ip)

	// Validate IP input
	validIp := net.ParseIP(ip)
	if validIp == nil {
		return nil, errors.New("Provided IP " + ip + " is not a valid IP.")
	}

	// Retrieve the lookup result from the database
	result, err := db.GetIPLookupResult(r.Database, validIp)
	if err != nil {
		log.Printf("error while retrieving lookup result: %s", err)
		return nil, err
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
