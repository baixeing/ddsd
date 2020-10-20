package ls

import (
	"net/http"

	"io/ioutil"

	"fmt"
	"os"
	"text/tabwriter"

	"encoding/json"

	"time"

	"github.com/baixeing/ddsd/cmd/ddsctl/resolver"
	"github.com/baixeing/ddsd/storage"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func Table(jsonData []byte) error {
	files := make(storage.Files, 0)

	if err := json.Unmarshal(jsonData, &files); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 2, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "UID\tFILENAME\tPATH\tSIZE\tCONTENT-TYPE")

	for _, f := range files {
		_, _ = fmt.Fprintf(w,
			"%s\t%s\t%s\t%s\t%s\n",
			f.UID,
			f.Name,
			f.Path,
			humanize.IBytes(f.Size),
			f.ContentType,
		)
	}
	return w.Flush()
}

var (
	Cmd = &cobra.Command{
		Use:                   "ls",
		Aliases:               []string{"list", "dir"},
		Short:                 "list files",
		DisableFlagsInUseLine: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			endpoint, err := resolver.Endpoint(time.Second)
			if err != nil {
				return err
			}

			endpoint.Path = "ls"
			r, err := http.Get(endpoint.String())
			if err != nil {
				return err
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}

			if err = Table(body); err != nil {
				return err
			}

			// var pretty bytes.Buffer
			// if err = json.Indent(&pretty, body, "", "  "); err != nil {
			// 	return err
			// }
			// cmd.Println(pretty.String())

			return r.Body.Close()
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)
