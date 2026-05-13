package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Board struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name" validate:"required,min=1,max=100"`
	Description string    `json:"description" validate:"max=500"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type List struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Card struct {
	ID          string    `json:"id"`
	ListID      string    `json:"list_id" validate:"required"`
	Title       string    `json:"title" validate:"required,min=1,max=200"`
	Description string    `json:"description" validate:"max=1000"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Board) Validate() error {
	return validate.Struct(b)
}

func (l *List) Validate() error {
	return validate.Struct(l)
}

func (c *Card) Validate() error {
	return validate.Struct(c)
}

func (u *User) Validate() error {
	return validate.Struct(u)
}
