BUILD = `date +%FT%T%z`
# APP_VERSION := `git describe --tags --abbrev=4`

_CURDIR := `git rev-parse --show-toplevel 2>/dev/null | sed -e 's/(//'`

PKG_LIST := $(shell go list ${_CURDIR}/... 2>/dev/null)

EXAMPLE := ${_CURDIR}/_example

linter := golangci-lint



.PHONY: all
all: lint

.PHONY: build_all
build_all: lint

.PHONY: full_tests
full_tests: go_lint_max go_performance unit_tests race msan ## Full tests

.PHONY: go_security
go_security: ## Check bugs
	@${linter} run --disable-all -E gosec -E govet \
	-E scopelint -E staticcheck -E typecheck

.PHONY: lint
lint: go_lint go_performance go_security ## Full liner checks

.PHONY: go_lint
go_lint: ## Lint the files
	@${linter} run

.PHONY: go_mod
go_mod: ## mod with proxy
	@GOPROXY="https://proxy.golang.org" \
	go mod download ${PKG_LIST}

.PHONY: go_performance
go_performance: ## Check performance
	@${linter} run --disable-all -p performance

.PHONY: go_lint_max
go_lint_max: ## Max lint checks the files
	@${linter} run -p bugs -p complexity -p unused -p performance -p format \
	-E interfacer -E gocritic

.PHONY: go_style
go_style: ## check style of code
	@${linter} run -p style

.PHONY: unit_tests
test: unit_tests ## Run unittests
unit_tests: ## Run unittests
	@go test -short ${PKG_LIST}

.PHONY: race
race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

.PHONY: bench
bench: ## Run benchmark tests
	@go test -benchmem -bench=. -run=^$ ${PKG_LIST}


.PHONY: msan
msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

.PHONY: coverage
coverage: ## Generate global code coverage report
	[ -x /opt/tools/bin/coverage.sh ] && /opt/tools/bin/coverage.sh || ${_CURDIR}/scripts/coverage.sh;

.PHONY: coverhtml
coverhtml: ## Generate global code coverage report in HTML
	[ -x /opt/tools/bin/coverage.sh ] && /opt/tools/bin/coverage.sh html || ${_CURDIR}/scripts/coverage.sh html;

.PHONY: dep
dep: ## Get the dependencies
	@go mod vendor



.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
