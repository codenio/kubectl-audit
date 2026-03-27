
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