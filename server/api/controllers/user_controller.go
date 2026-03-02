package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/api/dto"
	"server/api/service"
	"server/internal/middleware"
	"server/internal/response"
	"time"
)

type UserController interface {
	Create(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Refresh(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	Session(http.ResponseWriter, *http.Request)
}

type userController struct {
	userService        *service.UserService
	refreshCookieName  string
	secureRefreshToken bool
	refreshTTL         time.Duration
	refreshSameSite    http.SameSite
}

func NewUserController(
	us *service.UserService,
	refreshCookieName string,
	secureRefreshToken bool,
	refreshTTL time.Duration,
	refreshSameSite http.SameSite,
) UserController {
	if refreshCookieName == "" {
		refreshCookieName = "refresh_token"
	}
	return &userController{
		userService:        us,
		refreshCookieName:  refreshCookieName,
		secureRefreshToken: secureRefreshToken,
		refreshTTL:         refreshTTL,
		refreshSameSite:    refreshSameSite,
	}
}

func (uc userController) Create(w http.ResponseWriter, r *http.Request) {

	var p dto.UserRegistrationPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&p)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid request payload.")
		return
	}
	if fieldErrors := p.FieldErrors(); len(fieldErrors) > 0 {
		response.SendValidationError(w, http.StatusBadRequest, "Please correct the highlighted fields.", fieldErrors)
		return
	}

	_, err = uc.userService.Register(p)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPayload):
			response.SendError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid request payload.")
		case errors.Is(err, service.ErrUserExists):
			response.SendError(w, http.StatusConflict, "USER_ALREADY_EXISTS", "User already exists.")
		default:
			response.SendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error.")
		}
		return
	}

	response.SendMessage(w, "Account created!", http.StatusOK)
}

func (uc *userController) Login(w http.ResponseWriter, r *http.Request) {
	var p dto.UserLoginPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&p)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid request payload.")
		return
	}
	if fieldErrors := p.FieldErrors(); len(fieldErrors) > 0 {
		response.SendValidationError(w, http.StatusBadRequest, "Please correct the highlighted fields.", fieldErrors)
		return
	}

	user, err := uc.userService.Login(r.Context(), p, service.AuthClientMeta{
		IP:        r.RemoteAddr,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPayload):
			response.SendError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid request payload.")
		case errors.Is(err, service.ErrInvalidCredentials):
			response.SendError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials.")
		default:
			response.SendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error.")
		}
		return
	}

	uc.setRefreshCookie(w, user.RefreshToken)
	response.SendJSON(w, http.StatusOK, user)
}

func (uc *userController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(uc.refreshCookieName)
	if err != nil {
		response.SendError(w, http.StatusUnauthorized, "MISSING_REFRESH_TOKEN", "Missing refresh token.")
		return
	}

	result, err := uc.userService.Refresh(r.Context(), cookie.Value, service.AuthClientMeta{
		IP:        r.RemoteAddr,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingToken), errors.Is(err, service.ErrInvalidToken), errors.Is(err, service.ErrExpiredSession):
			response.SendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized.")
		case errors.Is(err, service.ErrTokenReuseDetected):
			response.SendError(w, http.StatusUnauthorized, "TOKEN_REUSE_DETECTED", "Session invalidated. Please login again.")
		default:
			response.SendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error.")
		}
		return
	}

	uc.setRefreshCookie(w, result.RefreshToken)
	response.SendJSON(w, http.StatusOK, result)
}

func (uc *userController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(uc.refreshCookieName)
	if err == nil {
		_ = uc.userService.Logout(r.Context(), cookie.Value)
	}

	uc.clearRefreshCookie(w)
	response.SendMessage(w, "logged out", http.StatusOK)
}

func (uc *userController) Session(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.SendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized.")
		return
	}

	session, err := uc.userService.Session(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			response.SendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized.")
		default:
			response.SendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error.")
		}
		return
	}

	response.SendJSON(w, http.StatusOK, session)
}

func (uc *userController) setRefreshCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     uc.refreshCookieName,
		Value:    refreshToken,
		Path:     "/auth",
		MaxAge:   int(uc.refreshTTL.Seconds()),
		HttpOnly: true,
		Secure:   uc.secureRefreshToken,
		SameSite: uc.refreshSameSite,
	})
}

func (uc *userController) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     uc.refreshCookieName,
		Value:    "",
		Path:     "/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   uc.secureRefreshToken,
		SameSite: uc.refreshSameSite,
	})
}
