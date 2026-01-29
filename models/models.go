package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Custom validator instance
var validate = validator.New()

// Board - доска
type Board struct {
	ID          string    `json:"id" validate:"omitempty,uuid"`
	Name        string    `json:"name" validate:"required,min=1,max=100"`
	Description string    `json:"description" validate:"max=500"`
	CreatedAt   time.Time `json:"created_at" validate:"omitempty"`
	UpdatedAt   time.Time `json:"updated_at" validate:"omitempty"`
}

// List - список на доске
type List struct {
	ID        string    `json:"id" validate:"omitempty,uuid"`
	BoardID   string    `json:"board_id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
	Position  int       `json:"position" validate:"min=0"`
	CreatedAt time.Time `json:"created_at" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" validate:"omitempty"`
}

// Card - карточка в списке
type Card struct {
	ID          string    `json:"id" validate:"omitempty,uuid"`
	ListID      string    `json:"list_id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required,min=1,max=200"`
	Description string    `json:"description" validate:"max=1000"`
	Position    int       `json:"position" validate:"min=0"`
	CreatedAt   time.Time `json:"created_at" validate:"omitempty"`
	UpdatedAt   time.Time `json:"updated_at" validate:"omitempty"`
}

// Validate validates the struct
func (b *Board) Validate() error {
	return validate.Struct(b)
}

func (l *List) Validate() error {
	return validate.Struct(l)
}

func (c *Card) Validate() error {
	return validate.Struct(c)
}
