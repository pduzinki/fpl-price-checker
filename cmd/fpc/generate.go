package fpc

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate-report",
	Aliases: []string{"gen"},
	Short:   "generates today's price changes report", // TODO later it should be possible to generate report for a date provided by the user
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fpc generate-report command")

		// TODO add later
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
