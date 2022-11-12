# vim: set foldmarker={,} foldlevel=0 foldmethod=marker:
#
#
# This Makefile is heavily inspired by:
# https://github.com/vincentbernat/hellogopher/blob/master/Makefile
#

# can set different values locally by adding a '.env' file
-include ./.env

export

.PHONY: help
help: ## Shows help message.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: deps
deps: ## Downloads dependencies.
	@echo 'Installing dependencies...'
	@go mod download

.PHONY: clean
clean: ## Cleans all the possible dirty files generated by the other go commands.
	@echo 'Checking for files created by go/coverage...'
	@[ -f ./profile.cov ] && rm profile.cov	|| true
	@[ -f ./profile.cov.tmp ] && rm profile.cov.tmp || true
	@echo 'checking for the bin folder...'
	@[ -d ./bin ] && rm -r ./bin || true


.PHONY: proto
proto: ## Generates the proto files according to the `/api` definitions.
	@[ -d ./pkg/api ] && rm -r ./pkg/api || true
	@echo 'Generating proto files...'
	@protoc api/*.proto --go_out=pkg --go_opt=paths=source_relative --proto_path=.

.PHONY: python
python: ## TODO.


