package main

import (
	"github.com/codenio/kubectl-audit/cmd/plugin/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // load auth providers for client-go
)

func main() {
	cli.InitAndExecute()
}
