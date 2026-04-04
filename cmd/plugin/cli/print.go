package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"
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

// WriteAuditSummary prints total (resources in scope), benign (healthy / desired by kind), and attention
// (rows displayed — audit matches) before kubectl-formatted output.
// resourceTitle is a singular label (e.g. "Node", "Pod"). When dest is a TTY, benign is green and attention is red.
func WriteAuditSummary(dest io.Writer, resourceTitle string, total, benign, attention int) {
	colorize := false
	if f, okf := dest.(*os.File); okf {
		colorize = term.IsTerminal(int(f.Fd()))
	}
	if colorize {
		g := color.New(color.FgGreen).SprintFunc()
		r := color.New(color.FgRed).SprintFunc()
		fmt.Fprintln(dest, "-------------------------------------------------------")
		fmt.Fprintf(dest, "%s Audit summary: total = %d benign = %s attention = %s\n", resourceTitle, total, g(benign), r(attention))
		fmt.Fprintln(dest, "-------------------------------------------------------")
		fmt.Fprintf(dest, "%s that requires attention\n", resourceTitle)
		return
	}
	fmt.Fprintln(dest, "-------------------------------------------------------")
	fmt.Fprintf(dest, "%s Audit summary: total = %d benign = %d attention = %d\n\n", resourceTitle, total, benign, attention)
	fmt.Fprintln(dest, "-------------------------------------------------------")
	fmt.Fprintf(dest, "%s that requires attention\n", resourceTitle)
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
