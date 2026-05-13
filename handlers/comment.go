package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexandervashurin/trello-golang/models"
	"github.com/alexandervashurin/trello-golang/storage"
	"github.com/alexandervashurin/trello-golang/utils"
	"github.com/google/uuid"
)

type CommentHandler struct {
	storage *storage.Storage
}

func NewCommentHandler(s *storage.Storage) *CommentHandler {
	return &CommentHandler{storage: s}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var req struct {
		CardID  string `json:"card_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.CardID == "" || req.Content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "card_id and content are required")
		return
	}

	boardID, err := h.storage.GetCardBoardID(req.CardID)
	if err != nil || boardID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Card not found")
		return
	}

	board, err := h.storage.GetBoard(boardID)
	if err != nil || board == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	if !board.IsPublic && board.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	now := time.Now()
	comment := &models.Comment{
		ID:        uuid.New().String(),
		CardID:    req.CardID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.storage.CreateComment(comment); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	userInfo, _ := h.storage.GetUserInfo(userID)
	if userInfo != nil {
		comment.Username = userInfo.Username
	}

	utils.RespondWithSuccess(w, comment)
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	cardID := r.URL.Query().Get("card_id")
	if cardID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "card_id is required")
		return
	}

	boardID, err := h.storage.GetCardBoardID(cardID)
	if err != nil || boardID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Card not found")
		return
	}

	board, err := h.storage.GetBoard(boardID)
	if err != nil || board == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	userID, _ := r.Context().Value(UserIDKey).(string)
	if !board.IsPublic && board.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	limit, offset := parsePagination(r)
	comments, err := h.storage.GetCommentsByCard(cardID, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondWithSuccess(w, comments)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	comment, err := h.storage.GetComment(id)
	if err != nil || comment == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	if comment.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.storage.DeleteComment(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	utils.RespondWithSuccess(w, map[string]string{"message": "Comment deleted successfully"})
}
