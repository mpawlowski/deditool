package cmd

import (
	"fmt"
	"os"

	"github.com/mpawlowski/deditool/v2/cmd/steam"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "deditool",
	Short: "Manage and troubleshoot dedicated gameservers",
	Long:  `Manage and troubleshoot dedicated gameservers`,
}


func init() {
	rootCmd.AddCommand(steam.Cmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}