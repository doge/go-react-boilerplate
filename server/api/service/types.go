package service

import "time"

type LoginResult struct {
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"-"`
}

type RefreshResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"-"`
}

type SessionResult struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
}

type AuthClientMeta struct {
	IP        string
	UserAgent string
}

type AuthSettings struct {
	RefreshTTL        time.Duration
	RefreshCookieName string
}
