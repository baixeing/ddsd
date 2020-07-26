package rm

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:                   "rm",
	Aliases:               []string{"remove", "rm", "del"},
	Short:                 "remove files",
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
