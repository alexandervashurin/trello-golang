package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexandervashurin/trello-golang/models"
	"github.com/alexandervashurin/trello-golang/storage"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	board.ID = uuid.New().String()
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	h.storage.CreateBoard(&board)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	board, exists := h.storage.GetBoard(id)
	if !exists {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func (h *Handler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
	boards := h.storage.GetAllBoards()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	h.storage.DeleteBoard(id)
	w.WriteHeader(http.StatusOK)
}

// List handlers
func (h *Handler) CreateList(w http.ResponseWriter, r *http.Request) {
	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list.ID = uuid.New().String()
	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()

	h.storage.CreateList(&list)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) GetListsByBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		http.Error(w, "board_id is required", http.StatusBadRequest)
		return
	}

	lists := h.storage.GetListsByBoard(boardID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lists)
}

func (h *Handler) DeleteList(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	h.storage.DeleteList(id)
	w.WriteHeader(http.StatusOK)
}

// Card handlers
func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card.ID = uuid.New().String()
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	h.storage.CreateCard(&card)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func (h *Handler) GetCardsByList(w http.ResponseWriter, r *http.Request) {
	listID := r.URL.Query().Get("list_id")
	if listID == "" {
		http.Error(w, "list_id is required", http.StatusBadRequest)
		return
	}

	cards := h.storage.GetCardsByList(listID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	h.storage.DeleteCard(id)
	w.WriteHeader(http.StatusOK)
}
