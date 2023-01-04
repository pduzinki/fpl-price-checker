package fpc

import (
	"fmt"

	"github.com/pduzinki/fpl-price-checker/pkg/di"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "start-server",
	Aliases: []string{"server"},
	Short:   "starts fpc web server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fpc start-server command")

		s := di.NewServer()

		s.Start(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
