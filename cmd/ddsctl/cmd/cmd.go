package cmd

import (
	"os"

	"github.com/baixeing/ddsd/cmd/ddsctl/cmd/completion"
	"github.com/baixeing/ddsd/cmd/ddsctl/cmd/get"
	"github.com/baixeing/ddsd/cmd/ddsctl/cmd/ls"
	"github.com/baixeing/ddsd/cmd/ddsctl/cmd/put"
	"github.com/baixeing/ddsd/cmd/ddsctl/cmd/rm"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:                   "ddsctl",
		Short:                 "ddsctl tool",
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
)

func init() {
	Cmd.SetHelpTemplate(Template)
	Cmd.AddCommand(ls.Cmd, rm.Cmd, put.Cmd, get.Cmd, completion.Cmd) // [TODO] subommands
}

func Execute() {
	if err := Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
