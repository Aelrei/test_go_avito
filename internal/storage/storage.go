package storage

import (
	"errors"
)

var (
	ErrTagNotFound     = errors.New("tag not found")
	ErrFeatureNotFound = errors.New("feature not found")
)
