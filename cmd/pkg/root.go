package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:               "mast",
		Short:             "mast control interface.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long: `mast configuration command line utility for service operators to debug and diagnose their istio.
`,
		//PersistentPreRunE: configureLogging,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(deployCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
