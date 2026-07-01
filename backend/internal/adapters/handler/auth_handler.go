package handler

import (
	"encoding/json"
	"errors"
	"go-hexagonal/internal/api/dto"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/ports"
	"net/http"
)

type AuthHandler struct {
	userService ports.UserService
}

func NewAuthHandler(service ports.UserService) *AuthHandler {
	return &AuthHandler{
		userService: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "Invalid request body-register")
		return
	}

	if err := validate.Struct(req); err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.userService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrorUserAlreadyExists) {
			h.WriteError(w, http.StatusConflict, err.Error())
			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully!"})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "Invalid request body-login")
	}

	if err := validate.Struct(req); err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
	}

	token, err := h.userService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrorInvalidCredentials) {
			h.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.AuthResponse{Token: token})
}

func (h *AuthHandler) WriteError(w http.ResponseWriter, status int, msg string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
