GO_GET := go get -u
GO_CLEAN := go clean
GO_TEST := go test
GO_TOOL := go tool

.PHONY: deps
deps:
	$(GO_GET) github.com/stretchr/testify/assert

.PHONY: clean
clean:
	$(GO_CLEAN)

.PHONY: test
test:
		$(GO_TEST) -v -short -cover github.com/davegarred/event-ledger...

.PHONY: coverage
coverage:
		$(GO_TEST) -short -coverprofile=.coverage.out github.com/davegarred/event-ledger...
		$(GO_TOOL) cover -html=.coverage.out