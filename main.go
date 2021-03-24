package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Result struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ResponseCode net.IP    `json:"response_code"`
	IPAddress    net.IP    `json:"ip_address"`
}

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

// LookupIP looks up the targetIP against the dnsblAddress
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

func main() {
	// Initialize database
	database, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err.Error())
	}
	defer database.Close()

	// Reverse the IP
	myIp := net.ParseIP("1.2.3.4")
	reversedIp, err := ReverseIP(myIp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Perform the lookup
	fmt.Println(reversedIp)
	response, err := LookupIP(reversedIp, "zen.spamhaus.org")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result := &Result{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ResponseCode: response,
		IPAddress:    reversedIp,
	}

	j, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("result: \n" + string(j))

	fmt.Println("storing in database")
	sqlStatement, err := database.Prepare("CREATE TABLE IF NOT EXISTS address_results (uuid TEXT UNIQUE PRIMARY KEY, created_at TEXT, updated_at TEXT, response_code TEXT, ip_address TEXT)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlStatement.Exec()

	insertStatement, err := database.Prepare("INSERT INTO address_results (uuid, created_at, updated_at, response_code, ip_address) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = insertStatement.Exec("1", result.CreatedAt.String(), result.UpdatedAt.String(), result.ResponseCode.String(), result.IPAddress.String())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rows, err := database.Query("SELECT response_code FROM address_results")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var rCode string
	for rows.Next() {
		rows.Scan(&rCode)
	}

	queriedResult := &Result{
		ResponseCode: net.ParseIP(rCode),
	}
	fmt.Println(rCode)
	fmt.Println("queried result: " + queriedResult.ResponseCode.String())
}
