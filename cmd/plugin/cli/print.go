package cli

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/cmd/get"
	kubectldscheme "k8s.io/kubectl/pkg/scheme"
)

// NewAuditPrintFlags returns kubectl get's PrintFlags (-o json/yaml/wide, custom-columns, go-template, etc.).
func NewAuditPrintFlags() *get.PrintFlags {
	return get.NewGetPrintFlags()
}

// printObjects prints obj using kubectl get printer wiring: TypeSetter + optional TablePrinter for server tables.
// get.PrintFlags.Copy() is shallow; HumanReadableFlags is shared, so WithNamespace is reset each call.
func printObjects(obj runtime.Object, pf *get.PrintFlags, gk schema.GroupKind, withNamespace, humanReadable bool) error {
	pf.HumanReadableFlags.WithNamespace = false
	pf.SetKind(gk)
	if withNamespace {
		_ = pf.EnsureWithNamespace()
	}

	printer, err := pf.ToPrinter()
	if err != nil {
		return err
	}
	printer, err = printers.NewTypeSetter(kubectldscheme.Scheme).WrapToPrinter(printer, nil)
	if err != nil {
		return err
	}
	if humanReadable {
		printer = &get.TablePrinter{Delegate: printer}
	}
	return printer.PrintObj(obj, os.Stdout)
}
