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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fpc start-server command")

		ctx := cmd.Context()
		s := di.NewServer()

		go func() {
			<-ctx.Done()
			os.Exit(0)
		}()

		s.Start(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
