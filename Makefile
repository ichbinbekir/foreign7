GO := go

run:
	@$(GO) run $(CURDIR)/cmd/foreign7

build:
	@$(GO) build -o $(CURDIR)/bin/ $(CURDIR)/cmd/foreign7
