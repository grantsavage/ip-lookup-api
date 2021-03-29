# GraphQL IP DNSBL Lookup API
__ip-lookup-api__ is a GraphQL service that queries and stores the Spamhaus Blocklist for malicious IP addresses.

## Building
### Docker
To bulid the service as a Docker container, run the following `docker` command inside the project directory:
```
docker build -t iplookup:1.0 .
```
### Executable
To build the service as an executable, make sure `go` version `1.16` is insalled and run:
```
make build
```
This will build the service into a `server` executable.

## Running
### Configuration
This service will respect the following environment variables:
|Variable Name|Description|Default|
|---|---|---|
|PORT|The port on which to bind the server to.|8080|
|AUTH_USERNAME|The username which requests will be authenticated against.||
|AUTH_PASSWORD|The password which requests will be authenticated against.||

### Docker
To run the service as a `docker` container and configure the necessary environment variables, use the following command:
```bash
docker run -p 3000:3000 -e PORT="3000" -e AUTH_USERNAME="secureworks" AUTH_PASSWORD="supersecret" iplookup:1.0
```
Optionally add the `-it` flag to view the container logs.

### Executable
To run the executable, use the following command:
```bash
PORT="3000" AUTH_USERNAME="secureworks" AUTH_PASSWORD="supersecret" ./server
```

## How to Use
### GraphQL Endpoint
The GraphQL endpoint for this service is at `/graphql`. 

### Authorization Token
First, create a basic authorization token by running the following command:
```bash
printf "%s:%s" "secureworks" "supersecret" | base64
```
Now use the generated token and set it as the `Authorization` header like so:
```
Authorization: Basic <your token here>
```

### Enqueue
With the authorization token set, you can enqueue IP addresses using by executing the following mutation at `/graphql`:
```graphql
mutation {
    enqueue(ips: ["1.2.3.4"])
}
```

### Get IP Details
With the authorization token set, you can query the lookup details of an IP by executing the following query:
```graphql
query {
    getIPDetails(ip: "1.2.3.4") {
        uuid
        ip_address
        response_code
        created_at
        updated_at
    }
}
```