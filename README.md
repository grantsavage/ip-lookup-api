# GraphQL IP DNSBL Lookup API
__ip-lookup-api__ is a GraphQL service that queries and stores the Spamhaus Blocklist for malicious IP addresses.

### Table of Contents  
* [Building](#building)  
* [Running](#running)  
* [How to Use](#how-to-use)
* [Project Structure](#project-structure)
* [Packages Used](#packages-used)

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

## Project Structure
I did my best to separate the core concerns of the application into 4 major packages: `auth`,`db`,`graph`, and `dns`.

* `auth` : Provides the authentication mechanism and HTTP middlware for the application.
* `db` : Provides an interface for the application to interact with the database. 
* `graph` : Defines and implements the resolvers for the GraphQL interface.
* `dns` : Provides methods to handle validating IP addresses and performing the DNS host lookup of an IP.

## Packages Used
* [99designs/gqlgen](https://github.com/99designs/gqlgen): Used to implement the GraphQL interface.
* [go-chi.chi](https://github.com/go-chi/chi) : Used to attach the authentiation middleware to the server.
* [mattn/go-sqlite3](http://github.com/mattn/go-sqlite3) : Used to interact with the SQLite database.
* [satori/go.uuid](https://github.com/satori/go.uuid) : Used to generate UUIDs for new IP lookup results.
* [vektah/gqlparser/v2](https://github.com/vektah/gqlparser/v2) : Used in conjunction with `99designs/gqlgen`.
* [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) : Use in `db` test suite.