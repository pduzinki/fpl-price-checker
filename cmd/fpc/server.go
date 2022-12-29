package fpc

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "start-server",
	Aliases: []string{"server"},
	Short:   "starts fpc web server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fpc start-server command")

		// TODO add later

		// TODO remove that later
		ch := make(chan int)
		<-ch // wait forever
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
