.PHONY: build build-linux run test clean deploy

BINARY  := profilepage
PORT    := 8080
# Date of the last commit that changed what the page shows; rendered in the footer.
UPDATED := $(shell git log -1 --format=%cs -- content templates static 2>/dev/null)
LDFLAGS := -X main.lastUpdated=$(UPDATED)

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w $(LDFLAGS)" -o $(BINARY) .

run:
	go run -ldflags="$(LDFLAGS)" .

test:
	go test ./...

clean:
	rm -f $(BINARY) $(BINARY)-linux-amd64

# Password comes from the DEPLOY_PASSWORD env var (falls back to SSH key auth when unset).
deploy:
	bash scripts/deploy.sh
