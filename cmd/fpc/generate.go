package fpc

import (
	"fmt"

	"github.com/pduzinki/fpl-price-checker/pkg/di"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate-report",
	Aliases: []string{"gen"},
	Short:   "generates today's price changes report", // TODO later it should be possible to generate report for a date provided by the user
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		gs, err := di.NewGenerateService()
		if err != nil {
			return fmt.Errorf("generate-report cmd failed: %w", err)
		}

		err = gs.GeneratePriceReport(ctx)
		if err != nil {
			return fmt.Errorf("generate-report cmd failed: %w", err)
		}

		fmt.Println("generate-report cmd finished")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
