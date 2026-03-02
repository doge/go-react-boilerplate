package service

import "errors"

var (
	ErrInvalidCredentials = errors.New(`invalid credentials`)
	ErrInternal           = errors.New(`internal error`)
	ErrInvalidToken       = errors.New(`invalid token`)
	ErrMissingToken       = errors.New(`missing token`)
	ErrExpiredSession     = errors.New(`expired session`)
	ErrTokenReuseDetected = errors.New(`refresh token reuse detected`)
	ErrInvalidPayload     = errors.New(`invalid payload`)
	ErrUserExists         = errors.New(`user already exists`)
)
