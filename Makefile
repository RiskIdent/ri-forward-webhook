.PHONY: build
build: ri-forward-webhook

ri-forward-webhook: *.go
	go build

.PHONY: clean
clean:
	rm -f ri-forward-webhook

.PHONY: deps
deps: deps-npm deps-go

.PHONY: deps-go
deps-go:
	go mod download
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: deps-npm
deps-npm: node_modules

node_modules: package.json
	npm install

.PHONY: lint
lint: lint-md lint-go

.PHONY: lint-fix
lint-fix: lint-md-fix lint-go-fix

.PHONY: lint-md
lint-md: node_modules
	npx markdownlint-cli2

.PHONY: lint-md-fix
lint-md-fix: node_modules
	npx markdownlint-cli2 --fix

.PHONY: lint-go
lint-go:
	@echo goimports -d '**/*.go'
	@goimports -d $(shell git ls-files "*.go")

.PHONY: lint-go-fix
lint-go-fix:
	@echo goimports -d -w '**/*.go'
	@goimports -d -w $(shell git ls-files "*.go")
