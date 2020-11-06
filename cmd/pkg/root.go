package cmd

import (
	"fmt"
	"github.com/mj37yhyy/mast/pkg/kubernetes"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:               "mast",
		Short:             "mast control interface.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long: `mast configuration command line utility for service operators to
debug and diagnose their istio.
`,
		Run: func(cmd *cobra.Command, args []string) {
			kubernetes.InitKubernetesClient(cfgFile)
		},
		//PersistentPreRunE: configureLogging,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"kubernetes config file (default is $HOME/.kube/config)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(deployCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
