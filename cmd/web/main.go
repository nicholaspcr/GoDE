package main

import (
	"log"

	"github.com/nicholaspcr/GoDE/cmd/web/internal/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
