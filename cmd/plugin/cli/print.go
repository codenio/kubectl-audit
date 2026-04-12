package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/fatih/color"
	"golang.org/x/term"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/cmd/get"
	kubectldscheme "k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/yaml"
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

// printContainersOutput prints a ContainerList: default/wide table, or json/yaml/name.
func printContainersOutput(list *plugin.ContainerList, out string, wide, allNS bool) error {
	if out == "" || out == "wide" {
		return printContainerTable(list, wide, allNS)
	}
	switch out {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(list)
	case "yaml":
		b, err := yaml.Marshal(list)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(b)
		return err
	case "name":
		for _, c := range list.Items {
			if c.Namespace != "" {
				fmt.Fprintf(os.Stdout, "%s/%s\n", c.Namespace, c.Name)
			} else {
				fmt.Fprintln(os.Stdout, c.Name)
			}
		}
		return nil
	default:
		return fmt.Errorf("containers: output format %q is not supported (use default table, wide, json, yaml, name)", out)
	}
}

func printContainerTable(list *plugin.ContainerList, wide, allNS bool) error {
	w := printers.GetNewTabWriter(os.Stdout)
	defer w.Flush()

	var hdr []string
	if allNS {
		hdr = append(hdr, "NAMESPACE")
	}
	hdr = append(hdr, "POD", "NAME")
	if wide {
		hdr = append(hdr, "READY", "STATUS", "RESTARTS", "AGE", "PORTS", "IMAGE", "PULLPOLICY", "TYPE")
	} else {
		hdr = append(hdr, "READY", "STATUS", "RESTARTS", "AGE", "TYPE")
	}
	fmt.Fprintln(w, strings.Join(hdr, "\t"))
	for i := range list.Items {
		c := &list.Items[i]
		var row []string
		if allNS {
			row = append(row, c.Namespace)
		}
		row = append(row, c.Pod, c.Name)
		if wide {
			row = append(row, c.Ready, c.Status, c.Restarts, c.Age, c.Ports, c.Image, c.PullPolicy, c.Type)
		} else {
			row = append(row, c.Ready, c.Status, c.Restarts, c.Age, c.Type)
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return nil
}
