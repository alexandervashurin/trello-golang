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

// Board handlers
func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var board models.Board
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Валидация
	if err := board.Validate(); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	board.ID = uuid.New().String()
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	h.storage.CreateBoard(&board)

	utils.RespondWithSuccess(w, board)
}

func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	board, exists := h.storage.GetBoard(id)
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	utils.RespondWithSuccess(w, board)
}

func (h *Handler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
	boards := h.storage.GetAllBoards()
	utils.RespondWithSuccess(w, boards)
}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	_, exists := h.storage.GetBoard(id)
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	h.storage.DeleteBoard(id)
	utils.RespondWithSuccess(w, map[string]string{"message": "Board deleted successfully"})
}

// List handlers
func (h *Handler) CreateList(w http.ResponseWriter, r *http.Request) {
	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Валидация
	if err := list.Validate(); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	// Проверка существования доски
	_, exists := h.storage.GetBoard(list.BoardID)
	if !exists {
		utils.RespondWithError(w, http.StatusBadRequest, "Board not found")
		return
	}

	list.ID = uuid.New().String()
	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()

	h.storage.CreateList(&list)

	utils.RespondWithSuccess(w, list)
}

func (h *Handler) GetListsByBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "board_id is required")
		return
	}

	lists := h.storage.GetListsByBoard(boardID)
	utils.RespondWithSuccess(w, lists)
}

func (h *Handler) DeleteList(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	_, exists := h.storage.GetList(id)
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	h.storage.DeleteList(id)
	utils.RespondWithSuccess(w, map[string]string{"message": "List deleted successfully"})
}

// Card handlers
func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Валидация
	if err := card.Validate(); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	// Проверка существования списка
	_, exists := h.storage.GetList(card.ListID)
	if !exists {
		utils.RespondWithError(w, http.StatusBadRequest, "List not found")
		return
	}

	card.ID = uuid.New().String()
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	h.storage.CreateCard(&card)

	utils.RespondWithSuccess(w, card)
}

func (h *Handler) GetCardsByList(w http.ResponseWriter, r *http.Request) {
	listID := r.URL.Query().Get("list_id")
	if listID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "list_id is required")
		return
	}

	cards := h.storage.GetCardsByList(listID)
	utils.RespondWithSuccess(w, cards)
}

func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	_, exists := h.storage.GetCard(id)
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	h.storage.DeleteCard(id)
	utils.RespondWithSuccess(w, map[string]string{"message": "Card deleted successfully"})
}
