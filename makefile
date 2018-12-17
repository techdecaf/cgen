HELP_SPACING=15
PACKAGE_NAME=cgen
# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty --abbrev=0)
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

install: ## Install all the build and lint dependencies
	# go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/dep/cmd/dep
	# gometalinter --install --update
	@$(MAKE) dep

build:
	go build -ldflags "-X main.VERSION=$(VERSION)" -o .dist/$(PACKAGE_NAME) -v

.PHONY: test cover dep clean release version
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

dep: ## Run dep ensure and prune
	dep ensure

clean: ## Remove temporary files
	go clean

version: ## prints the current version tag
	@echo $(VERSION)

release:
	git push origin $(VERSION)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Print help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-$(HELP_SPACING)s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help