package model

import "errors"

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ErrConflict is returned when an operation conflicts with existing data (e.g. duplicate email).
var ErrConflict = errors.New("conflict")
