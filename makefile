HELP_SPACING=15
PACKAGE_NAME=cgen
# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

install: ## Install all the build and lint dependencies
	# go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/dep/cmd/dep
	# gometalinter --install --update
	@$(MAKE) dep

.PHONY: bump-patch
bump-patch:
	$(eval NEW_VERSION := $(shell echo $(VERSION)| awk -F. -v OFS=. -v f=3 '{ $$f++ } 1'))
	git tag -a $(NEW_VERSION) -m "patch release"
bump-minor:
	$(eval NEW_VERSION := $(shell echo $(VERSION)| awk -F. -v OFS=. -v f=2 '{ $$f++ } 1'))
	git tag -a $(NEW_VERSION) -m "minor release"
bump-major:
	$(eval NEW_VERSION := $(shell echo $(VERSION)| awk -F. -v OFS=. -v f=1 '{ $$f++ } 1'))
	git tag -a $(NEW_VERSION) -m "breaking change"

build: 
	go build -ldflags "-X main.VERSION=$(VERSION)" -o $(PACKAGE_NAME) -v

.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: dep
dep: ## Run dep ensure and prune
	dep ensure

.PHONY: clean
clean: ## Remove temporary files
	go clean

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Print help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-$(HELP_SPACING)s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help