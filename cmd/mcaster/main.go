package main

import (
	"log"

	"github.com/hyposcaler-bot/mcaster/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
