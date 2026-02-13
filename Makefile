BINARY = mdview
VERSION = 0.1.0
LDFLAGS = -ldflags "-s -w -X github.com/bilabl/mdview/cmd.Version=$(VERSION)"
BINDIR = bin
INSTALL_DIR = $(HOME)/.local/bin
SKILL_DIR = $(HOME)/.claude/skills/mdview

.PHONY: build install test clean generate css

build: css
	@mkdir -p $(BINDIR)
	go build $(LDFLAGS) -o $(BINDIR)/$(BINARY) .

install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BINDIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo ""
	@echo "Installed $(BINARY) to $(INSTALL_DIR)/$(BINARY)"
	@# Check if INSTALL_DIR is in PATH
	@echo "$$PATH" | tr ':' '\n' | grep -qx "$(INSTALL_DIR)" \
		|| echo "\n⚠️  $(INSTALL_DIR) is NOT in your PATH!\n   Add this to your shell profile (~/.bashrc or ~/.zshrc):\n\n   export PATH=\"\$$HOME/.local/bin:\$$PATH\"\n\n   Then reload: source ~/.bashrc"
	@# Copy SKILL.md for Claude Code integration
	@if [ ! -d "$(SKILL_DIR)" ]; then \
		mkdir -p "$(SKILL_DIR)"; \
		echo "Created $(SKILL_DIR)"; \
	fi
	@if [ -f SKILL.md ] && [ ! -f "$(SKILL_DIR)/SKILL.md" ]; then \
		cp SKILL.md "$(SKILL_DIR)/SKILL.md"; \
		echo "Copied SKILL.md to $(SKILL_DIR)/SKILL.md"; \
	fi

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
