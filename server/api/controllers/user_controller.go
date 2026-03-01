package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/api/dto"
	"server/api/service"
	"server/internal/response"
)

type UserController interface {
	Create(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type userController struct {
	userService *service.UserService
}

func NewUserController(us *service.UserService) UserController {
	return &userController{userService: us}
}

func (uc userController) Create(w http.ResponseWriter, r *http.Request) {

	var p *dto.UserRegistrationPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&p)
	if err != nil {
		response.SendMessage(w, "invalid payload", http.StatusBadRequest)
		return
	}

	_, err = uc.userService.Register(*p)
	if err != nil {
		response.SendMessage(w, err.Error(), http.StatusInternalServerError)
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
		response.SendMessage(w, "invalid payload", http.StatusBadRequest)
		return
	}

	user, err := uc.userService.Login(r.Context(), p)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.SendMessage(w, "invalid credentials", http.StatusUnauthorized)
		default:
			response.SendMessage(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	response.SendJSON(w, http.StatusOK, user)
}
