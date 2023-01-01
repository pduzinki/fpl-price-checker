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
		fetchService, err := di.NewFetchService()
		if err != nil {
			return fmt.Errorf("fetch-fpl-data cmd failed: %w", err)
		}

		err = fetchService.Fetch()
		if err != nil {
			return fmt.Errorf("fetch-fpl-data cmd failed: %w", err)
		}

		fmt.Println("fetch-fpl-data cmd finished")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
