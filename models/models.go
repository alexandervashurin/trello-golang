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

type Comment struct {
	ID        string    `json:"id"`
	CardID    string    `json:"card_id" validate:"required"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) Validate() error {
	return validate.Struct(c)
}

type Attachment struct {
	ID         string    `json:"id"`
	CardID     string    `json:"card_id" validate:"required"`
	UserID     string    `json:"user_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"-"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	CreatedAt  time.Time `json:"created_at"`
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
