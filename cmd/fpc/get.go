package fpc

import (
	"encoding/json"
	"fmt"

	"github.com/pduzinki/fpl-price-checker/pkg/di"
	"github.com/spf13/cobra"
)

var getReportCmd = &cobra.Command{
	Use:     "get-report",
	Aliases: []string{"get"},
	Short:   "gets today's price changes report",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		gs := di.NewGetService()

		report, err := gs.GetLatestReport(ctx)
		if err != nil {
			return fmt.Errorf("get-report cmd failed: %w", err)
		}

		jsonReport, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("get-report cmd failed: %w", err)
		}

		fmt.Println(string(jsonReport))

		fmt.Println("get-report cmd finished")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getReportCmd)
}
