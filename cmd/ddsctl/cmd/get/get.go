package get

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:                   "get",
	Aliases:               []string{"pull"},
	Short:                 "get files from DDSD",
	DisableFlagsInUseLine: true,

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
