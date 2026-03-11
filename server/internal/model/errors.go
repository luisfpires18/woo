package model

import "errors"

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ErrConflict is returned when an operation conflicts with existing data (e.g. duplicate email).
var ErrConflict = errors.New("conflict")

// ErrInsufficientResources is returned when a player lacks enough resources for an action.
var ErrInsufficientResources = errors.New("insufficient resources")

// ErrInsufficientGold is returned when a player lacks enough gold for an action.
var ErrInsufficientGold = errors.New("insufficient gold")

// ErrBuildingInProgress is returned when a construction queue slot is already occupied.
var ErrBuildingInProgress = errors.New("building already under construction")

// ErrMaxLevelReached is returned when a building is already at its maximum level.
var ErrMaxLevelReached = errors.New("building already at max level")

// ErrPrerequisitesNotMet is returned when building prerequisites are not satisfied.
var ErrPrerequisitesNotMet = errors.New("prerequisites not met")
