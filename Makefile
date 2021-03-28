build: lint
	go build -o server

run: build
	./server

lint:
	go fmt ./...
	go vet ./...

test: lint
	go test -v -coverprofile=coverage.out ./services ./db ./auth