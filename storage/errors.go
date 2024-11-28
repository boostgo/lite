package storage

import "errors"

var (
	// ErrConnNotSelected returns if "shard client" does not choose connection to use
	ErrConnNotSelected = errors.New("connection not selected")
)
