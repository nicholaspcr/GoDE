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
help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

### Go
.PHONY: go-download
go-download: ## download dependencies
	go mod download

.PHONY: clean
clean: ## cleans all the possible dirty files generated by the other go commands
	@echo 'checking for files created by go/coverage...'
	@[ -f ./profile.cov ] && rm profile.cov	|| true
	@[ -f ./profile.cov.tmp ] && rm profile.cov.tmp || true
	@echo 'checking for the bin folder...'
	@[ -d ./bin ] && rm -r ./bin || true


.PHONY: proto
proto: ## generates the proto files according to the `/api` definitions
	protoc api/v1/*.proto --go_out=pkg --go_opt=paths=source_relative --proto_path=.

### Python

.PHONY: python
python: ## TODO


