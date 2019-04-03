HELP_SPACING=15

CI_PROJECT_NAME ?= cgen
CI_COMMIT_TAG ?= $(shell git describe --tags --always --dirty --abbrev=0)
CI_COMMIT_REF_NAME ?= $(shell git rev-parse --abbrev-ref HEAD)  # curernt branch
CI_COMMIT_SHA ?= $(shell git rev-parse HEAD)

S3_BUCKET := github.techdecaf.io

install: ## Install all the build and lint dependencies
	# go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/dep/cmd/dep
	# gometalinter --install --update
	dep ensure

build: clean ## complie window, linux, osx x64
	# @export GOARCH=amd64
	GOOS=windows go build -ldflags "-X main.VERSION=$(CI_COMMIT_TAG)" -o .dist/windows/$(CI_PROJECT_NAME).exe -v
	GOOS=darwin go build -ldflags "-X main.VERSION=$(CI_COMMIT_TAG)" -o .dist/osx/$(CI_PROJECT_NAME) -v
	GOOS=linux go build -ldflags "-X main.VERSION=$(CI_COMMIT_TAG)" -o .dist/linux/$(CI_PROJECT_NAME) -v

test: ## Run all the tests
	@echo no tests yet, you should add some.

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

clean: ## Remove temporary files
	go clean
	rm -rf .dist

version: ## prints the current version tag
	@echo $(CI_COMMIT_TAG)

publish: ## push build to s3
	@aws s3 sync .dist s3://$(S3_BUCKET)/$(CI_PROJECT_NAME)/$(CI_COMMIT_TAG)
	@aws s3 sync .dist s3://$(S3_BUCKET)/$(CI_PROJECT_NAME)/latest

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-$(HELP_SPACING)s\033[0m %s\n", $$1, $$2}'

.PHONY: install build test cover dep clean version publish help
.DEFAULT_GOAL := help

