.PHONY: build build-linux run test clean

BINARY := profilepage
PORT   := 8080

build:
	go build -o $(BINARY) .

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) .

run:
	go run .

test:
	go test ./...

clean:
	rm -f $(BINARY) $(BINARY)-linux-amd64
