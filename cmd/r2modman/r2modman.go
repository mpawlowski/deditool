package r2modman

import (
	"fmt"

	"github.com/spf13/cobra"
)

var r2modmanCmd = &cobra.Command{
	Use:   "r2modman",
	Short: "Tools for r2ModmanPlus profile exports.",
	Long:  `Tools for r2ModmanPlus profile exports.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("specify a command")
	},
}

func init() {
	r2modmanCmd.AddCommand(syncCmd)
}

func Cmd() *cobra.Command {
	return r2modmanCmd
}
