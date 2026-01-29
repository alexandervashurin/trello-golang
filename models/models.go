package models

import "time"

// Board - доска
type Board struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// List - список на доске
type List struct {
    ID        string    `json:"id"`
    BoardID   string    `json:"board_id"`
    Name      string    `json:"name"`
    Position  int       `json:"position"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Card - карточка в списке
type Card struct {
    ID          string    `json:"id"`
    ListID      string    `json:"list_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Position    int       `json:"position"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
