BUILD = `date +%FT%T%z`
# APP_VERSION := `git describe --tags --abbrev=4`

_CURDIR := `git rev-parse --show-toplevel 2>/dev/null | sed -e 's/(//'`

PKG_LIST := $(shell go list ${_CURDIR}/... 2>/dev/null | grep -v internal/examples | grep -v internal/build)


linter := golangci-lint

mockery := mockery

test := go test


.PHONY: all
all: test lint

.PHONY: full_tests
full_tests: go_lint_max unit_tests race msan ## Full tests

.PHONY: lint
lint: go_lint go_security ## Full liner checks



.PHONY: go_security
go_security: ## Check bugs
	@${linter} run --disable-all \
	-E gosec -E govet -E exportloopref -E staticcheck -E typecheck

.PHONY: go_lint
go_lint: ## Lint the files
	@${linter} run

.PHONY: go_lint_max
go_lint_max: ## Max lint checks the files
	@${linter} run \
	-p bugs -p complexity -p unused -p format \
	-E gosec -E govet -E exportloopref -E staticcheck -E typecheck

.PHONY: go_style
go_style: ## check style of code
	@${linter} run -p style

.PHONY: go_mod
go_mod: ## mod with proxy
	@GOPROXY="https://proxy.golang.org" \
	go mod download ${PKG_LIST}

.PHONY: test
test: unit_test test_paltform race msan ## Run unittests

.PHONY: test_paltform
test_paltform: ##Run paltform test
	@${test} -tags $(${_CURDIR}/build/check) ${PKG_LIST}

.PHONY: unit_test
unit_test: ## Run unittests
	@${test} ${PKG_LIST}

.PHONY: race
race: ## Run data race detector
	@${test} -race ${PKG_LIST}

.PHONY: msan
msan: ## Run memory sanitizer
	@CXX=clang++ CC=clang \
	${test} -msan ${PKG_LIST}


.PHONY: bench
bench: ## Run benchmark tests
	@${test} -benchmem -run=^# -bench=. ${PKG_LIST}



.PHONY: coverage
coverage: ## Generate global code coverage report
	@go test -cover ${PKG_LIST}


.PHONY: coverhtml
coverhtml: ## Generate global code coverage report in HTML
	@[ -x /opt/tools/bin/coverage.sh ] && /opt/tools/bin/coverage.sh html || \
	${_CURDIR}/scripts/tools/coverage.sh html ${_CURDIR}/build/coverage.html


.PHONY: run
example := "monitor"
run: ## Run example, call: make example={fiber|monitor} run
	@go run ${_CURDIR}/internal/examples/cmd/${example}

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
