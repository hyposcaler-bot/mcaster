package main

import (
	"log"

	"github.com/yourusername/mcaster/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
