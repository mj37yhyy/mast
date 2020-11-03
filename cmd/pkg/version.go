package cmd

import "github.com/spf13/cobra"

var (
	version      = "1.0.0"
	istioVersion = "1.7.2"
	versionCmd   = &cobra.Command{
		Use:               "version",
		Short:             "version for mast.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("mast version %s support istio version %s.", version, istioVersion)
		},
	}
)
