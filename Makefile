.PHONY: fmt
fmt: ## Format source using gofmt
	@gofumpt -l -w .

imports: ## fix go imports
	@goimports -local github.com/wal1251/pkg -w -l ./
#	@gci write --skip-generated -s standard,default,"prefix(github.com/wal1251/pkg)" .

lint: ## Linter for golang
	@docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.56.2-alpine golangci-lint run ./...

lint-fix: ## Linter fixes for golang
	@docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.56.2-alpine golangci-lint run --fix ./...

test:  ## Run all unit test
	@go test ./... -race -cover -short -v

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)