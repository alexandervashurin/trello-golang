package storage

import (
    "sync"
    "github.com/yourusername/go-trello/models"
)

type Storage struct {
    boards map[string]*models.Board
    lists  map[string]*models.List
    cards  map[string]*models.Card

    boardsMu sync.RWMutex
    listsMu  sync.RWMutex
    cardsMu  sync.RWMutex
}

func NewStorage() *Storage {
    return &Storage{
        boards: make(map[string]*models.Board),
        lists:  make(map[string]*models.List),
        cards:  make(map[string]*models.Card),
    }
}

// Boards
func (s *Storage) CreateBoard(board *models.Board) {
    s.boardsMu.Lock()
    defer s.boardsMu.Unlock()
    s.boards[board.ID] = board
}

func (s *Storage) GetBoard(id string) (*models.Board, bool) {
    s.boardsMu.RLock()
    defer s.boardsMu.RUnlock()
    board, exists := s.boards[id]
    return board, exists
}

func (s *Storage) GetAllBoards() []*models.Board {
    s.boardsMu.RLock()
    defer s.boardsMu.RUnlock()

    boards := make([]*models.Board, 0, len(s.boards))
    for _, board := range s.boards {
        boards = append(boards, board)
    }
    return boards
}

func (s *Storage) DeleteBoard(id string) {
    s.boardsMu.Lock()
    defer s.boardsMu.Unlock()
    delete(s.boards, id)
}

// Lists
func (s *Storage) CreateList(list *models.List) {
    s.listsMu.Lock()
    defer s.listsMu.Unlock()
    s.lists[list.ID] = list
}

func (s *Storage) GetList(id string) (*models.List, bool) {
    s.listsMu.RLock()
    defer s.listsMu.RUnlock()
    list, exists := s.lists[id]
    return list, exists
}

func (s *Storage) GetListsByBoard(boardID string) []*models.List {
    s.listsMu.RLock()
    defer s.listsMu.RUnlock()

    lists := make([]*models.List, 0)
    for _, list := range s.lists {
        if list.BoardID == boardID {
            lists = append(lists, list)
        }
    }
    return lists
}

func (s *Storage) DeleteList(id string) {
    s.listsMu.Lock()
    defer s.listsMu.Unlock()
    delete(s.lists, id)
}

// Cards
func (s *Storage) CreateCard(card *models.Card) {
    s.cardsMu.Lock()
    defer s.cardsMu.Unlock()
    s.cards[card.ID] = card
}

func (s *Storage) GetCard(id string) (*models.Card, bool) {
    s.cardsMu.RLock()
    defer s.cardsMu.RUnlock()
    card, exists := s.cards[id]
    return card, exists
}

func (s *Storage) GetCardsByList(listID string) []*models.Card {
    s.cardsMu.RLock()
    defer s.cardsMu.RUnlock()

    cards := make([]*models.Card, 0)
    for _, card := range s.cards {
        if card.ListID == listID {
            cards = append(cards, card)
        }
    }
    return cards
}

func (s *Storage) DeleteCard(id string) {
    s.cardsMu.Lock()
    defer s.cardsMu.Unlock()
    delete(s.cards, id)
}
