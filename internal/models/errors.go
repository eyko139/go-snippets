package models

import (
	"errors"
)

// Introduce custom errors to further decouple handlers from data store specific types
var (
  ErrNoRecord = errors.New("models: No matching record found")
  ErrInvalidCredentials = errors.New("models: invalid credentials")
  ErrDuplicateEmail = errors.New("models: duplicate email")
)
