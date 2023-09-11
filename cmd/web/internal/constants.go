// Package internal contains internal constants for the package.
package internal

import (
	"os"
	"path"
	"strings"
)

const (
	// projectSuffix is the default path for the project.
	projectSuffix = "cmd/web"
)

// ProjectPath returns the base path for the `web` project folder, its root is
// found on the GoDE/cmd/web folder. The purpose of this is to have a consistent
// path to relative files even when running the binary from a different
// directory localy.
func ProjectPath() string {
	currDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(currDir, projectSuffix) {
		return currDir
	}
	return path.Join(currDir, projectSuffix)
}
