package fpc

import (
	"fmt"

	"github.com/pduzinki/fpl-price-checker/pkg/di"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:     "fetch-fpl-data",
	Aliases: []string{"fetch"},
	Short:   "fetches current player prices data from FPL API",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("fpc fetch-fpl-data command")

		fetchService := di.NewFetchService()

		err := fetchService.Fetch()
		if err != nil {
			return fmt.Errorf("fetch-fpl-data cmd failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
