PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

.PHONY: wordle install uninstall test

wordle:
	go build -o wordle .

test:
	go test ./...

install: wordle
	mkdir -p "$(BINDIR)"
	cp wordle "$(BINDIR)/wordle"
	@echo "Installed $(BINDIR)/wordle"
	@echo "If 'wordle' is not found, add this to your shell rc:"
	@echo "  export PATH=\"\$$PATH:$(BINDIR)\""

uninstall:
	rm -f "$(BINDIR)/wordle"
