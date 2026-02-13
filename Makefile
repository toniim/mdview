BINARY = mdview
VERSION = 0.1.0
LDFLAGS = -ldflags "-s -w -X github.com/bilabl/mdview/cmd.Version=$(VERSION)"
BINDIR = bin

.PHONY: build install test clean generate css

build: css
	@mkdir -p $(BINDIR)
	go build $(LDFLAGS) -o $(BINDIR)/$(BINARY) .

install: build
	cp $(BINDIR)/$(BINARY) $(HOME)/.local/bin/$(BINARY)

css:
	@if [ ! -f assets/chroma-light.css ] || [ ! -f assets/chroma-dark.css ]; then \
		go run tools/generate_chroma_css.go; \
	fi

generate:
	go run tools/generate_chroma_css.go

test:
	go test ./...

clean:
	rm -rf $(BINDIR)
	rm -f assets/chroma-light.css assets/chroma-dark.css

# Cross-compile targets
.PHONY: build-all
build-all: css
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe .
