GO      := $(shell which go)
GOFUMPT := $(shell which gofumpt)
LINT    := $(shell which golangci-lint)

ifeq ($(GO),)
$(error Go not found in PATH)
endif

.PHONY: test lint fmt fmt-check clean

test:
	$(GO) test ./...

lint:
ifeq ($(LINT),)
	$(error golangci-lint not found in PATH)
endif
	$(LINT) run

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

clean:
	$(GO) clean ./...
