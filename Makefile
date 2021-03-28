build: lint
	go build -o server

generate: go generate ./...

lint:
	go fmt ./...
	go vet ./...

run: build
	./server

test: lint
	go test -v -coverprofile=coverage.out ./services ./db ./auth