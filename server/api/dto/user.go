package dto

type UserRegistrationPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
