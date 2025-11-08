// Package main implements the DE server providing gRPC and HTTP APIs for DE optimization.
package main

import (
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/commands"
)

func main() { commands.Execute() }
