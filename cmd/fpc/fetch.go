package fpc

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:     "fetch-fpl-data",
	Aliases: []string{"fetch"},
	Short:   "fetches current player prices data from FPL API",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fpc fetch-fpl-data command")

		// TODO add later
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
