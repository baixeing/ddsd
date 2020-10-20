package put

import (
	"log"
	"time"

	"os"

	"net/http"

	"net/http/httputil"

	"fmt"

	"github.com/baixeing/ddsd/cmd/ddsctl/resolver"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:                   "put [FILES...]",
	Aliases:               []string{"push"},
	Short:                 "put files or directories to DDSD",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),

	PreRunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, err := resolver.Endpoint(time.Second)
		if err != nil {
			return err
		}
		endpoint.Path = "put"

		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				log.Println(err)
				continue
			}
			r, err := http.Post(endpoint.String(), "application/octet-stream", f)
			if err != nil {
				log.Println(err)
				f.Close()
				continue
			}
			b, err := httputil.DumpResponse(r, true)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(b))

			r.Body.Close()
			f.Close()

		}

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
