package main

import (
	"log"

	cmd "github.com/kgf1980/go-luxpower/cmd/go-luxpower"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
