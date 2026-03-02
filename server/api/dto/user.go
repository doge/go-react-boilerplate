package dto

import (
	"net/mail"
	"regexp"
	"strings"
)

type UserRegistrationPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

type FieldErrors map[string]string

func (e FieldErrors) Error() string {
	return "validation failed"
}

func (p UserRegistrationPayload) Validate() error {
	if fieldErrors := p.FieldErrors(); len(fieldErrors) > 0 {
		return fieldErrors
	}

	return nil
}

func (p UserRegistrationPayload) FieldErrors() FieldErrors {
	errs := FieldErrors{}

	p.Email = strings.TrimSpace(p.Email)
	p.Username = strings.TrimSpace(p.Username)

	if p.Email == "" {
		errs["email"] = "Email is required."
	}
	if p.Username == "" {
		errs["username"] = "Username is required."
	}
	if p.Password == "" {
		errs["password"] = "Password is required."
	}

	if p.Email != "" {
		if _, err := mail.ParseAddress(p.Email); err != nil {
			errs["email"] = "Email format is invalid."
		}
	}
	if p.Username != "" && !usernameRegex.MatchString(p.Username) {
		errs["username"] = "Username must be 3-32 chars and only use letters, numbers, _ or -."
	}
	if p.Password != "" && (len(p.Password) < 8 || len(p.Password) > 72) {
		errs["password"] = "Password must be between 8 and 72 characters."
	}

	return errs
}

func (p UserLoginPayload) Validate() error {
	if fieldErrors := p.FieldErrors(); len(fieldErrors) > 0 {
		return fieldErrors
	}

	return nil
}

func (p UserLoginPayload) FieldErrors() FieldErrors {
	errs := FieldErrors{}

	p.Username = strings.TrimSpace(p.Username)
	if p.Username == "" {
		errs["username"] = "Username is required."
	}
	if p.Password == "" {
		errs["password"] = "Password is required."
	}
	if p.Username != "" && len(p.Username) < 3 {
		errs["username"] = "Username must be at least 3 characters."
	}
	if p.Password != "" && len(p.Password) > 72 {
		errs["password"] = "Password must be 72 characters or less."
	}

	return errs
}
