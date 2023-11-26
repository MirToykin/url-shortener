package storage

import "errors"

var (
	ErrAliasExists = errors.New("alias exists")
	ErrURLNotFound = errors.New("URL not found")
)
