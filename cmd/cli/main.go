package main

import (
	"log"

	"github.com/mamcer/nostalgia/cmd/cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil && err.Error() != "" {
		log.Fatal(err)
	}
}
