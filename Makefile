
export GO111MODULE=on
# If proxy.golang.org times out or resets (common on VPN), Go falls back to direct source fetches.
export GOPROXY ?= https://proxy.golang.org,direct

GO_PACKAGES := ./pkg/... ./cmd/...
COVERPKG := github.com/codenio/kubectl-audit/...
GOTESTSUM ?= gotestsum
JUNIT_FILE := junit.xml
COVER_FILE := cover.out
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/codenio/kubectl-audit/pkg/version.Version=$(VERSION) \
	-X github.com/codenio/kubectl-audit/pkg/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/codenio/kubectl-audit/pkg/version.BuildDate=$(BUILD_DATE)

.PHONY: test
test:
	go test $(GO_PACKAGES) -coverprofile $(COVER_FILE) -coverpkg=$(COVERPKG)

.PHONY: test-ci
test-ci:
	@command -v $(GOTESTSUM) >/dev/null 2>&1 || { echo "$(GOTESTSUM) not found; run: go install gotest.tools/gotestsum@v1.13.0"; exit 1; }
	$(GOTESTSUM) --junitfile $(JUNIT_FILE) --format standard-verbose -- \
		-coverprofile=$(COVER_FILE) -covermode=atomic -coverpkg=$(COVERPKG) \
		$(GO_PACKAGES)

.PHONY: cover
cover: test
	go tool cover -func $(COVER_FILE)

.PHONY: bin
bin: fmt vet
	go build -ldflags "$(LDFLAGS)" -o bin/audit github.com/codenio/kubectl-audit/cmd/plugin

.PHONY: fmt
fmt:
	go fmt $(GO_PACKAGES)

.PHONY: vet
vet:
	go vet $(GO_PACKAGES)

.PHONY: precommit
precommit:
	@command -v pre-commit >/dev/null 2>&1 || { echo "pre-commit not found; install with: pip install pre-commit (or brew install pre-commit)"; exit 1; }
	pre-commit run go-fmt --all-files
	pre-commit run go-vet --all-files

.PHONY: precommit-install
precommit-install:
	@command -v pre-commit >/dev/null 2>&1 || { echo "pre-commit not found; install with: pip install pre-commit (or brew install pre-commit)"; exit 1; }
	pre-commit install

.PHONY: ci
ci: precommit test-ci bin

.PHONY: kubernetes-deps
kubernetes-deps:
	go get k8s.io/client-go@v0.28.4 k8s.io/api@v0.28.4 k8s.io/apimachinery@v0.28.4 \
		k8s.io/cli-runtime@v0.28.4 k8s.io/kubectl@v0.28.4

.PHONY: setup
setup:
	make -C setup

.PHONY: install
install: bin
	mkdir -p "$(HOME)/.krew/bin"
	cp bin/audit "$(HOME)/.krew/bin/kubectl-audit"
	chmod +x "$(HOME)/.krew/bin/kubectl-audit"
