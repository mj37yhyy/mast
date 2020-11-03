package cmd

import "github.com/spf13/cobra"

var (
	deployCmd = &cobra.Command{
		Use:               "deploy",
		Short:             "Deploy the service to kubernetes and use istio to publish in grayscale.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || len(args) > 1 {
				cmd.Help()
				return
			}
		},
	}
)

func init() {
	deployCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config)")
}
