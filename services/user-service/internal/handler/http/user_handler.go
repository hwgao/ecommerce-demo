package http

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "ecommerce/services/user-service/internal/domain"
    "ecommerce/shared/pkg/middleware"
    "ecommerce/shared/pkg/response"
)

type UserHandler struct {
    userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) SetupRoutes(router *mux.Router) {
    router.HandleFunc("/register", h.Register).Methods("POST")
    router.HandleFunc("/login", h.Login).Methods("POST")
    
    // Protected routes
    protected := router.PathPrefix("/profile").Subrouter()
    protected.Use(middleware.JWTAuth)
    protected.HandleFunc("", h.GetProfile).Methods("GET")
    protected.HandleFunc("", h.UpdateProfile).Methods("PUT")
}

type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    user, err := h.userService.Register(req.Email, req.Password, req.FirstName, req.LastName)
    if err != nil {
        response.Error(w, http.StatusConflict, err.Error())
        return
    }

    response.Success(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    token, err := h.userService.Login(req.Email, req.Password)
    if err != nil {
        response.Error(w, http.StatusUnauthorized, err.Error())
        return
    }

    response.Success(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(uuid.UUID)
    
    user, err := h.userService.GetProfile(userID)
    if err != nil {
        response.Error(w, http.StatusNotFound, "User not found")
        return
    }

    response.Success(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(uuid.UUID)
    
    var req struct {
        FirstName string `json:"first_name" validate:"required"`
        LastName  string `json:"last_name" validate:"required"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := h.userService.UpdateProfile(userID, req.FirstName, req.LastName); err != nil {
        response.Error(w, http.StatusInternalServerError, "Failed to update profile")
        return
    }

    response.Success(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
