HELP_SPACING=15
PACKAGE_NAME=cgen
PACKAGE_VERSION := $(shell git describe --tags --always --dirty --abbrev=0)
# This version-strategy uses git tags to set the version string
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

install: ## Install all the build and lint dependencies
	# go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/dep/cmd/dep
	# gometalinter --install --update
	@$(MAKE) dep

build: clean ## complie window, linux, osx x64
	# @export GOARCH=amd64
	GOOS=windows go build -ldflags "-X main.VERSION=$(PACKAGE_VERSION)" -o .dist/windows/$(PACKAGE_NAME).exe -v
	GOOS=darwin go build -ldflags "-X main.VERSION=$(PACKAGE_VERSION)" -o .dist/osx/$(PACKAGE_NAME) -v

test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

dep: ## Run dep ensure and prune
	dep ensure

clean: ## Remove temporary files
	go clean
	rm -rf .dist

version: ## prints the current version tag
	@echo $(PACKAGE_VERSION)

publish: ## push build to s3
	aws s3 sync ./dist $(S3_BUCKET)/$(PACKAGE_NAME)/$(PACKAGE_VERSION)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-$(HELP_SPACING)s\033[0m %s\n", $$1, $$2}'

.PHONY: install build test cover dep clean version publish help
.DEFAULT_GOAL := help

