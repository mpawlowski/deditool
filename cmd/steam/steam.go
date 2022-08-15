package steam

import (
	"fmt"

	"github.com/spf13/cobra"
)

var steamCmd = &cobra.Command{
	Use:   "steam",
	Short: "Tools to interact with Steam",
	Long:  `Tools to interact with Steam`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("specify a command")
	},
}

func init() {
	steamCmd.AddCommand(queryCmd, simCmd)
}

func Cmd() *cobra.Command {
	return steamCmd
}
