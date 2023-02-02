package fpc

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fpc",
	Short: "fpc - fpl price checker",
	Long:  `fpc is a simple app for checking player price changes in Fantasy Premier League`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s'\n", err)
		os.Exit(1)
	}
}
