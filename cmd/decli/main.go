package main

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
