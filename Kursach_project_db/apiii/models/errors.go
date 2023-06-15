package models

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("item not found")
	ErrConflict            = errors.New("item with such id already exists")
	ErrBadRequest          = errors.New("parashniy zapros")
	ErrInternalServerError = errors.New("src server error")
	ERRUnauthorized        = errors.New("no such user")
)
