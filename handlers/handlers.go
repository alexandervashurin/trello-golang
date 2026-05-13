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

type Handler struct {
	storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var board models.Board
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	if err := board.Validate(); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	board.ID = uuid.New().String()
	board.UserID = userID
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	if err := h.storage.CreateBoard(&board); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create board")
		return
	}

	utils.RespondWithSuccess(w, board)
}

func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	board, err := h.storage.GetBoard(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if board == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	if board.IsPublic {
		utils.RespondWithSuccess(w, board)
		return
	}

	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID != board.UserID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	utils.RespondWithSuccess(w, board)
}

func (h *Handler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	boards, err := h.storage.GetAllBoards(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondWithSuccess(w, boards)
}

func (h *Handler) GetPublicBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.storage.GetAllPublicBoards()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	utils.RespondWithSuccess(w, boards)
}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	ownerID, err := h.storage.GetBoardOwner(id)
	if err != nil || ownerID == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}
	if ownerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.storage.DeleteBoard(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete board")
		return
	}

	utils.RespondWithSuccess(w, map[string]string{"message": "Board deleted successfully"})
}

func (h *Handler) CreateList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	if err := list.Validate(); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	ownerID, err := h.storage.GetBoardOwner(list.BoardID)
	if err != nil || ownerID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Board not found")
		return
	}
	if ownerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	list.ID = uuid.New().String()
	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()

	if err := h.storage.CreateList(&list); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create list")
		return
	}

	utils.RespondWithSuccess(w, list)
}

func (h *Handler) GetListsByBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "board_id is required")
		return
	}

	board, err := h.storage.GetBoard(boardID)
	if err != nil || board == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	userID, ok := r.Context().Value(UserIDKey).(string)
	if !board.IsPublic && (!ok || userID != board.UserID) {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	lists, err := h.storage.GetListsByBoard(boardID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondWithSuccess(w, lists)
}

func (h *Handler) DeleteList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	list, err := h.storage.GetList(id)
	if err != nil || list == nil {
		utils.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	ownerID, err := h.storage.GetBoardOwner(list.BoardID)
	if err != nil || ownerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.storage.DeleteList(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete list")
		return
	}

	utils.RespondWithSuccess(w, map[string]string{"message": "List deleted successfully"})
}

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	if err := card.Validate(); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	list, err := h.storage.GetList(card.ListID)
	if err != nil || list == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "List not found")
		return
	}

	ownerID, err := h.storage.GetBoardOwner(list.BoardID)
	if err != nil || ownerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	card.ID = uuid.New().String()
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	if err := h.storage.CreateCard(&card); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create card")
		return
	}

	utils.RespondWithSuccess(w, card)
}

func (h *Handler) GetCardsByList(w http.ResponseWriter, r *http.Request) {
	listID := r.URL.Query().Get("list_id")
	if listID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "list_id is required")
		return
	}

	list, err := h.storage.GetList(listID)
	if err != nil || list == nil {
		utils.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	board, err := h.storage.GetBoard(list.BoardID)
	if err != nil || board == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	userID, ok := r.Context().Value(UserIDKey).(string)
	if !board.IsPublic && (!ok || userID != board.UserID) {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	cards, err := h.storage.GetCardsByList(listID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondWithSuccess(w, cards)
}

func (h *Handler) MoveCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var req struct {
		ID       string `json:"id"`
		ListID   string `json:"list_id"`
		Position int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	if req.ID == "" || req.ListID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "id and list_id are required")
		return
	}

	card, err := h.storage.GetCard(req.ID)
	if err != nil || card == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	list, err := h.storage.GetList(req.ListID)
	if err != nil || list == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Target list not found")
		return
	}

	ownerID, err := h.storage.GetBoardOwner(list.BoardID)
	if err != nil || ownerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	card.ListID = req.ListID
	card.Position = req.Position
	card.UpdatedAt = time.Now()

	if err := h.storage.UpdateCard(card); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to move card")
		return
	}

	utils.RespondWithSuccess(w, card)
}

func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	card, err := h.storage.GetCard(id)
	if err != nil || card == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	list, err := h.storage.GetList(card.ListID)
	if err != nil || list == nil {
		utils.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	ownerID, err := h.storage.GetBoardOwner(list.BoardID)
	if err != nil || ownerID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.storage.DeleteCard(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete card")
		return
	}

	utils.RespondWithSuccess(w, map[string]string{"message": "Card deleted successfully"})
}
