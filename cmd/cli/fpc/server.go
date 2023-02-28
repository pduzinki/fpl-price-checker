package fpc

import (
	"fmt"
	"os"

	"github.com/pduzinki/fpl-price-checker/pkg/di"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "start-server",
	Aliases: []string{"server"},
	Short:   "starts fpc web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("fpc start-server command")

		ctx := cmd.Context()
		s := di.NewServer()

		go func() {
			<-ctx.Done()
			os.Exit(0)
		}()

		if err := s.Start(":8080"); err != nil {
			return fmt.Errorf("start-server cmd failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
