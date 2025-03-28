VERSION ?= 0.0.0-dev

.PHONY: clean
clean:
	rm -rf release/
	rm -rf dist/

.PHONY: lint
lint: ## Run golang-ci-lint against code.
	go tool golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golang-ci-lint against code and apply fixes.
	go tool golangci-lint run --fix ./...

## Location to create the release
RELEASE_DIR ?= $(shell pwd)/release
$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)

.PHONY: release
release: $(RELEASE_DIR) ## Run release artifacts
	CGO_ENABLED=0 go tool gox -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="linux darwin windows" -arch="amd64 arm64"
	VERSION=$(VERSION) go tool gomplate -f hack/provider.yaml.tpl > $(RELEASE_DIR)/provider.yaml
	mv dist/* $(RELEASE_DIR)
	rm -rf dist/