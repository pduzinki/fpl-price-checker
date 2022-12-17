package main

import (
	"fmt"

	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"
)

func main() {
	fmt.Println("hello there, from fpl-price-checker")

	wr := wrapper.NewWrapper()

	player, err := wr.GetPlayers()
	if err != nil {
		panic(err)
	}

	_ = player

	fmt.Println(player[423])
}
