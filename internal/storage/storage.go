// internal/storage/storage.go

package storage

import "errors"

var (
	ErrDataNotFound = errors.New("data not found")
	ErrURLExists    = errors.New("data exists")
)
