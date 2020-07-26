package put

import (
	"time"

	"github.com/baixeing/ddsd/cmd/ddsctl/resolver"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:                   "put",
	Aliases:               []string{"push"},
	Short:                 "put files or directories to DDSD",
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, err := resolver.Endpoint(time.Second)
		if err != nil {
			return err
		}

		endpoint.Path = "put"

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
