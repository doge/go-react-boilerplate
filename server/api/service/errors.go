package service

import "errors"

var (
	ErrInvalidCredentials = errors.New(`invalid credentials`)
	ErrInternal           = errors.New(`internal error`)
)
