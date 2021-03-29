# Builds executable
build:
	go build -o server

# Builds docker container
container:
	docker build -t iplookup:$(tag) .

# Generates GraphQL code
generate: go generate ./...

# Formats and statically checks code
lint:
	go fmt ./...
	go vet ./...

# Builds and runs the application
run: build
	./server

# Runs the application test suites
test: lint
	go test -v -covermode=count -coverprofile=coverage.out ./dns ./db ./auth