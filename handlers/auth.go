package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexandervashurin/trello-golang/models"
	"github.com/alexandervashurin/trello-golang/storage"
	"github.com/alexandervashurin/trello-golang/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	storage *storage.Storage
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string          `json:"token"`
	User  *storage.UserInfo `json:"user"`
}

func NewAuthHandler(s *storage.Storage) *AuthHandler {
	return &AuthHandler{storage: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	if len(req.Password) < 6 {
		utils.RespondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	existing, _ := h.storage.GetUserByEmail(req.Email)
	if existing != nil {
		utils.RespondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	now := time.Now()
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.storage.CreateUser(user); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	userInfo := &storage.UserInfo{ID: user.ID, Username: user.Username, Email: user.Email}
	utils.RespondWithSuccess(w, AuthResponse{Token: token, User: userInfo})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.storage.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	userInfo := &storage.UserInfo{ID: user.ID, Username: user.Username, Email: user.Email}
	utils.RespondWithSuccess(w, AuthResponse{Token: token, User: userInfo})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	userInfo, err := h.storage.GetUserInfo(userID)
	if err != nil || userInfo == nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	utils.RespondWithSuccess(w, userInfo)
}

func (h *AuthHandler) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
