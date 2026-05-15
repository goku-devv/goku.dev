.PHONY: build run test clean

BINARY := goku.dev
PORT   := 8080

build:
	go build -o $(BINARY) .

run:
	go run .

test:
	go test ./...

clean:
	rm -f $(BINARY)
