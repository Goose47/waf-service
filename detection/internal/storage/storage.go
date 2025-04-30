// Package storage contains common storage items, such as errors.
package storage

import "errors"

// ErrNotFound represents common not found entity error.
var ErrNotFound = errors.New("not found")
