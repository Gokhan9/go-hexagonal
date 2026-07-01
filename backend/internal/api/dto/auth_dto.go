package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required", min=3, max50"`
	Password string `json:"password" validate:"required", min=6, max=100"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
