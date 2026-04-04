package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	// AuditPrintFlags is kubectl get's PrintFlags (k8s.io/kubectl/pkg/cmd/get).
	AuditPrintFlags *get.PrintFlags
)

func RootCmd() *cobra.Command {
	AuditPrintFlags = NewAuditPrintFlags()

	cmd := &cobra.Command{
		Use:   "audit [resource]",
		Short: "Run cluster audits with kubectl-compatible output",
		Long: `List Kubernetes resources that fail common health checks, using the same
output formats as kubectl get (-o wide, json, yaml, custom-columns, go-template, etc.).

For default and wide output, a one-line audit summary (total, benign, attention) is printed
before the resource table; see the README for sample output per resource.`,
		Example: `  kubectl audit pods
  kubectl audit pods -o wide
  kubectl audit nodes -o json
  kubectl audit pvc -o yaml
  kubectl audit jobs -o custom-columns=NAME:.metadata.name
  kubectl audit pods -A --selector app=nginx`,
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := normalizeResource(args[0])
			if err != nil {
				return err
			}
			allNS, _ := cmd.PersistentFlags().GetBool("all-namespaces")
			if !allNS {
				allNS = allNamespacesInArgv()
			}
			sel, _ := cmd.Flags().GetString("selector")

			opts := plugin.AuditOptions{
				AllNamespaces: allNS,
				LabelSelector: sel,
			}

			gk, ok := auditGroupKinds[res]
			if !ok {
				return fmt.Errorf("unknown resource %q", args[0])
			}

			var (
				obj           runtime.Object
				totalInScope  int
				benignInScope int
			)
			switch res {
			case "pods":
				obj, totalInScope, benignInScope, err = plugin.AuditPods(KubernetesConfigFlags, opts)
			case "nodes":
				obj, totalInScope, benignInScope, err = plugin.AuditNodes(KubernetesConfigFlags, opts)
			case "pv":
				obj, totalInScope, benignInScope, err = plugin.AuditPV(KubernetesConfigFlags, opts)
			case "pvc":
				obj, totalInScope, benignInScope, err = plugin.AuditPVC(KubernetesConfigFlags, opts)
			case "jobs":
				obj, totalInScope, benignInScope, err = plugin.AuditJobs(KubernetesConfigFlags, opts)
			case "cronjobs":
				obj, totalInScope, benignInScope, err = plugin.AuditCronJobs(KubernetesConfigFlags, opts)
			default:
				return fmt.Errorf("unknown resource %q", res)
			}
			if err != nil {
				return err
			}

			out := ""
			if AuditPrintFlags.OutputFormat != nil {
				out = *AuditPrintFlags.OutputFormat
			}
			if AuditPrintFlags.TemplateFlags != nil &&
				AuditPrintFlags.TemplateFlags.TemplateArgument != nil &&
				len(*AuditPrintFlags.TemplateFlags.TemplateArgument) > 0 && out == "" {
				out = "go-template"
			}

			if isHumanTableOutput(out) && auditObjectLen(obj) == 0 {
				total, benign, attention := plugin.SummarizeAudit(totalInScope, benignInScope, obj)
				WriteAuditSummary(os.Stdout, auditSummaryResourceTitle(res), total, benign, attention)
				writeAuditEmptyMessage(os.Stderr, res, allNS, KubernetesConfigFlags, totalInScope)
				return nil
			}

			total, benign, attention := plugin.SummarizeAudit(totalInScope, benignInScope, obj)

			obj, err = plugin.AsServerTableIfNeeded(KubernetesConfigFlags, res, opts, obj, out)
			if err != nil {
				return err
			}

			summaryDest := os.Stderr
			if isHumanTableOutput(out) {
				summaryDest = os.Stdout
			}
			WriteAuditSummary(summaryDest, auditSummaryResourceTitle(res), total, benign, attention)

			withNS := withNamespaceColumn(res, allNS)
			// Server Tables already include a NAMESPACE column when listing cluster-wide. Setting
			// PrintOptions.WithNamespace makes decorateTable prepend another column from row metadata;
			// row.Object is often unset (only Raw), so that prepends an empty cell and shifts columns.
			if _, ok := obj.(*metav1.Table); ok {
				withNS = false
			}
			return printObjects(obj, AuditPrintFlags, gk, withNS, isHumanTableOutput(out))
		},
	}

	cmd.PersistentFlags().BoolP("all-namespaces", "A", false, "If true, check the specified resource across all namespaces")
	cmd.Flags().StringP("selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='")

	AuditPrintFlags.AddFlags(cmd)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	KubernetesConfigFlags.AddFlags(cmd.PersistentFlags())

	return cmd
}

var auditGroupKinds = map[string]schema.GroupKind{
	"pods":     {Group: "", Kind: "Pod"},
	"nodes":    {Group: "", Kind: "Node"},
	"pv":       {Group: "", Kind: "PersistentVolume"},
	"pvc":      {Group: "", Kind: "PersistentVolumeClaim"},
	"jobs":     {Group: "batch", Kind: "Job"},
	"cronjobs": {Group: "batch", Kind: "CronJob"},
}

func auditSummaryResourceTitle(resource string) string {
	if gk, ok := auditGroupKinds[resource]; ok {
		return gk.Kind
	}
	return "Resource"
}

// allNamespacesInArgv catches -A / --all-namespaces if pflag/cobra did not (e.g. flag order with
// kubectl get PrintFlags or PersistentFlags merges).
func allNamespacesInArgv() bool {
	for _, a := range os.Args[1:] {
		switch a {
		case "-A", "--all-namespaces":
			return true
		}
		if strings.HasPrefix(a, "--all-namespaces=") {
			v := strings.TrimPrefix(a, "--all-namespaces=")
			return v == "" || v == "true" || v == "1" || v == "True"
		}
	}
	return false
}

func withNamespaceColumn(resource string, allNS bool) bool {
	if !allNS {
		return false
	}
	switch resource {
	case "pods", "pvc", "jobs", "cronjobs":
		return true
	default:
		return false
	}
}

func isHumanTableOutput(out string) bool {
	return out == "" || out == "wide"
}

func auditObjectLen(obj runtime.Object) int {
	if t, ok := obj.(*metav1.Table); ok {
		return len(t.Rows)
	}
	return meta.LenList(obj)
}

// writeAuditEmptyMessage prints a line after an empty audit table. inScopeCount is how many
// resources were listed before filtering (same basis as the summary total). When it is zero,
// nothing exists in scope; when positive, resources exist but none match the audit.
func writeAuditEmptyMessage(w io.Writer, resource string, allNS bool, cf *genericclioptions.ConfigFlags, inScopeCount int) {
	phrase := auditResourcePhrase(resource)
	namespaced := resource == "pods" || resource == "pvc" || resource == "jobs" || resource == "cronjobs"

	if inScopeCount > 0 {
		if namespaced && !allNS {
			ns, _, err := cf.ToRawKubeConfigLoader().Namespace()
			if err != nil || ns == "" {
				fmt.Fprintf(w, "No %s require attention.\n", phrase)
				return
			}
			fmt.Fprintf(w, "No %s require attention in %s namespace.\n", phrase, ns)
			return
		}
		fmt.Fprintf(w, "No %s require attention.\n", phrase)
		return
	}

	if namespaced && !allNS {
		ns, _, err := cf.ToRawKubeConfigLoader().Namespace()
		if err != nil || ns == "" {
			fmt.Fprintf(w, "No %s found.\n", phrase)
			return
		}
		fmt.Fprintf(w, "No %s found in %s namespace.\n", phrase, ns)
		return
	}
	fmt.Fprintf(w, "No %s found.\n", phrase)
}

func auditResourcePhrase(resource string) string {
	switch resource {
	case "pods":
		return "pods"
	case "nodes":
		return "nodes"
	case "pv":
		return "persistent volumes"
	case "pvc":
		return "persistent volume claims"
	case "jobs":
		return "jobs"
	case "cronjobs":
		return "cron jobs"
	default:
		return "resources"
	}
}

func normalizeResource(s string) (string, error) {
	switch strings.ToLower(s) {
	case "pod", "po", "pods":
		return "pods", nil
	case "node", "no", "nodes":
		return "nodes", nil
	case "pv", "pvs", "persistentvolume", "persistentvolumes":
		return "pv", nil
	case "pvc", "pvcs", "persistentvolumeclaim", "persistentvolumeclaims":
		return "pvc", nil
	case "job", "jobs":
		return "jobs", nil
	case "cronjob", "cronjobs", "cj":
		return "cronjobs", nil
	default:
		return "", fmt.Errorf("unsupported resource %q (supported: pods, nodes, pv, pvc, jobs, cronjobs)", s)
	}
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
