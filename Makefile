build:
	go build -o server

generate: go generate ./...

lint:
	go fmt ./...
	go vet ./...

run: build
	./server

test: lint
	go test -coverprofile=coverage.out ./dns ./db ./auth