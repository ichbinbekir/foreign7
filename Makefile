GO := go

run:
	@$(GO) run $(CURDIR)/cmd/foreign7

build:
	@mkdir -p $(CURDIR)/bin
	@cp -r $(CURDIR)/assets $(CURDIR)/bin/
	@$(GO) build -o $(CURDIR)/bin/foreign7 $(CURDIR)/cmd/foreign7

clean:
	@rm -rf $(CURDIR)/bin
