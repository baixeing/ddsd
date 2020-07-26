package completion

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:                   "completion",
		Short:                 "Completion code for the specified shell",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("bad")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				if err := cmd.Root().GenZshCompletion(cmd.OutOrStdout()); err != nil {
					return err
				}
				// [BUG] cobra ZSH completion
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), "compdef _ddsctl ddsctl"); err != nil {
					return err
				}
				return nil
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
			}

			return nil
		},
	}
)
