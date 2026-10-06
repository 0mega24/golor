GO      := $(shell which go)
GOFUMPT := $(shell which gofumpt)
LINT    := $(shell which golangci-lint)

ifeq ($(GO),)
$(error Go not found in PATH)
endif

.PHONY: test test-full bench fuzz vet vuln lint fix fmt fmt-check build cli clean

test:
	$(GO) test -race -count=1 ./...

# Run exhaustive conversion and CSS round-trip sweeps over all 16,777,216 sRGB colors.
test-full:
	GOLOR_FULL_SWEEP=1 $(GO) test -count=1 -timeout 30m ./...

bench:
	$(GO) test -run '^$$' -bench . -benchmem ./...

FUZZTIME ?= 30s
fuzz:
	$(GO) test -run '^$$' -fuzz FuzzParseCSS -fuzztime $(FUZZTIME) .
	$(GO) test -run '^$$' -fuzz FuzzHex -fuzztime $(FUZZTIME) .

vet:
	$(GO) vet ./...

vuln:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...

lint:
ifeq ($(LINT),)
	$(error golangci-lint not found in PATH)
endif
	$(LINT) run

fix:
ifeq ($(LINT),)
	$(error golangci-lint not found in PATH)
endif
	$(LINT) run --fix

fmt:
ifeq ($(GOFUMPT),)
	$(error gofumpt not found in PATH)
endif
	$(GOFUMPT) -l -w .

fmt-check:
ifeq ($(GOFUMPT),)
$(error gofumpt not found in PATH)
endif
	@if [ -n "$(shell $(GOFUMPT) -l .)" ]; then \
		echo "Files not gofumpt-formatted:"; \
		$(GOFUMPT) -l .; \
		exit 1; \
	fi

build:
	$(GO) build -o bin/golor ./cmd/golor

cli:
	@:

clean:
	$(GO) clean ./...
	rm -f bin/golor
