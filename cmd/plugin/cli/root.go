package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"strings"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Run cluster audits",
		Long:  "Run Kubernetes resource audits via subcommands.",
		Example: `  kubectl audit pods
  kubectl audit pods -n default
  kubectl audit nodes
  kubectl audit pvc -A
  kubectl audit jobs -A
  kubectl audit cronjobs -A
  kubectl audit pv`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newPodsCmd())
	cmd.AddCommand(newNodesCmd())
	cmd.AddCommand(newPVCCmd())
	cmd.AddCommand(newJobsCmd())
	cmd.AddCommand(newCronJobsCmd())
	cmd.AddCommand(newPVCmd())

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	KubernetesConfigFlags.AddFlags(cmd.PersistentFlags())

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
