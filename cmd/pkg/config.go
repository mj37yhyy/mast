package cmd

import (
	"github.com/mj37yhyy/mast/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:               "config",
		Short:             "config kubernetes config file for mast.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				if err := cmd.Help(); err != nil {
					cmd.PrintErr(err)
				}
				return
			}
			if err := kubernetes.InitKubernetesClient(args[0]); err != nil {
				cmd.PrintErr(err)
			}
			return
		},
	}
)
