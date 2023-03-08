
# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

clean_testcache:  ## Expire all Go test caches
	@echo "Cleaning test caches..."
	go clean -testcache ./...

unit_test:  ## Run all unit tests in the pkg folder
	@echo "Running unit test framework..."
	go test -v ./pkg/...

test:  unit_test ## Run package tests

lint_local: ## Run local instance of linting across the code base
	golangci-lint run ./...

lint_docker: ## Run docker instance of linting across the code base
	docker run --rm -v $(PWD):/app -v ~/.cache/golangci-lint/v1.50.1:/root/.cache -w /app golangci/golangci-lint:v1.50.1 golangci-lint run ./...
