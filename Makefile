
export GO111MODULE=on
# If proxy.golang.org times out or resets (common on VPN), Go falls back to direct source fetches.
export GOPROXY ?= https://proxy.golang.org,direct

.PHONY: test
test:
	go test ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: bin
bin: fmt vet
	go build -o bin/audit github.com/codenio/kubectl-audit/cmd/plugin

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: kubernetes-deps
kubernetes-deps:
	go get k8s.io/client-go@v11.0.0
	go get k8s.io/api@kubernetes-1.14.0
	go get k8s.io/apimachinery@kubernetes-1.14.0
	go get k8s.io/cli-runtime@kubernetes-1.14.0

.PHONY: setup
setup:
	make -C setup

.PHONY: install
install: bin
	mkdir -p "$(HOME)/.krew/bin"
	cp bin/audit "$(HOME)/.krew/bin/kubectl-audit"
	chmod +x "$(HOME)/.krew/bin/kubectl-audit"