// Package api provides the Golang implementation of the API defined in the
// proto files.
package api

import "google.golang.org/protobuf/types/known/emptypb"

// Empty is a helper variable to avoid creating a new empty message every time.
var Empty = &emptypb.Empty{}
