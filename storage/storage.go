package storage

import (
	"context"
	"fmt"

	"github.com/alexandervashurin/trello-golang/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func paginate(limit, offset int) string {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

// Boards
func (s *Storage) CreateBoard(board *models.Board) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO boards (id, user_id, name, description, is_public, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		board.ID, board.UserID, board.Name, board.Description, board.IsPublic, board.CreatedAt, board.UpdatedAt,
	)
	return err
}

func scanBoard(row pgx.Row) (*models.Board, error) {
	var b models.Board
	err := row.Scan(&b.ID, &b.UserID, &b.Name, &b.Description, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func scanBoards(rows pgx.Rows) ([]*models.Board, error) {
	defer rows.Close()
	var boards []*models.Board
	for rows.Next() {
		var b models.Board
		if err := rows.Scan(&b.ID, &b.UserID, &b.Name, &b.Description, &b.IsPublic, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, &b)
	}
	if boards == nil {
		boards = []*models.Board{}
	}
	return boards, nil
}

func (s *Storage) GetBoard(id string) (*models.Board, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, user_id, name, description, is_public, created_at, updated_at FROM boards WHERE id = $1`, id,
	)
	return scanBoard(row)
}

func (s *Storage) GetAllBoards(userID string, limit, offset int) ([]*models.Board, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, user_id, name, description, is_public, created_at, updated_at FROM boards WHERE user_id = $1 ORDER BY created_at DESC`+paginate(limit, offset),
		userID,
	)
	if err != nil {
		return nil, err
	}
	return scanBoards(rows)
}

func (s *Storage) GetAllPublicBoards(limit, offset int) ([]*models.Board, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, user_id, name, description, is_public, created_at, updated_at FROM boards WHERE is_public = TRUE ORDER BY created_at DESC`+paginate(limit, offset),
	)
	if err != nil {
		return nil, err
	}
	return scanBoards(rows)
}

func (s *Storage) DeleteBoard(id string) error {
	_, err := s.pool.Exec(context.Background(), `DELETE FROM boards WHERE id = $1`, id)
	return err
}

// Lists
func (s *Storage) CreateList(list *models.List) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO lists (id, board_id, name, position, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		list.ID, list.BoardID, list.Name, list.Position, list.CreatedAt, list.UpdatedAt,
	)
	return err
}

func (s *Storage) GetList(id string) (*models.List, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, board_id, name, position, created_at, updated_at FROM lists WHERE id = $1`, id,
	)
	var l models.List
	err := row.Scan(&l.ID, &l.BoardID, &l.Name, &l.Position, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (s *Storage) GetListsByBoard(boardID string) ([]*models.List, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, board_id, name, position, created_at, updated_at FROM lists WHERE board_id = $1 ORDER BY position ASC`,
		boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []*models.List
	for rows.Next() {
		var l models.List
		if err := rows.Scan(&l.ID, &l.BoardID, &l.Name, &l.Position, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, &l)
	}
	if lists == nil {
		lists = []*models.List{}
	}
	return lists, nil
}

func (s *Storage) DeleteList(id string) error {
	_, err := s.pool.Exec(context.Background(), `DELETE FROM lists WHERE id = $1`, id)
	return err
}

// Cards
func (s *Storage) CreateCard(card *models.Card) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO cards (id, list_id, title, description, position, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		card.ID, card.ListID, card.Title, card.Description, card.Position, card.CreatedAt, card.UpdatedAt,
	)
	return err
}

func (s *Storage) GetCard(id string) (*models.Card, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, list_id, title, description, position, created_at, updated_at FROM cards WHERE id = $1`, id,
	)
	var c models.Card
	err := row.Scan(&c.ID, &c.ListID, &c.Title, &c.Description, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *Storage) GetCardsByList(listID string) ([]*models.Card, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, list_id, title, description, position, created_at, updated_at FROM cards WHERE list_id = $1 ORDER BY position ASC`,
		listID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		var c models.Card
		if err := rows.Scan(&c.ID, &c.ListID, &c.Title, &c.Description, &c.Position, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, &c)
	}
	if cards == nil {
		cards = []*models.Card{}
	}
	return cards, nil
}

func (s *Storage) UpdateCard(card *models.Card) error {
	_, err := s.pool.Exec(context.Background(),
		`UPDATE cards SET list_id = $1, position = $2, updated_at = $3 WHERE id = $4`,
		card.ListID, card.Position, card.UpdatedAt, card.ID,
	)
	return err
}

func (s *Storage) DeleteCard(id string) error {
	_, err := s.pool.Exec(context.Background(), `DELETE FROM cards WHERE id = $1`, id)
	return err
}

// Comments
func (s *Storage) CreateComment(c *models.Comment) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO comments (id, card_id, user_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.CardID, c.UserID, c.Content, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (s *Storage) GetCommentsByCard(cardID string, limit, offset int) ([]models.Comment, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT c.id, c.card_id, c.user_id, u.username, c.content, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.card_id = $1
		ORDER BY c.created_at ASC`+paginate(limit, offset),
		cardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.CardID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []models.Comment{}
	}
	return comments, nil
}

func (s *Storage) GetComment(id string) (*models.Comment, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT c.id, c.card_id, c.user_id, u.username, c.content, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.id = $1`, id,
	)
	var c models.Comment
	err := row.Scan(&c.ID, &c.CardID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *Storage) GetCardBoardID(cardID string) (string, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT l.board_id FROM cards c JOIN lists l ON l.id = c.list_id WHERE c.id = $1`, cardID,
	)
	var boardID string
	err := row.Scan(&boardID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return boardID, nil
}

func (s *Storage) DeleteComment(id string) error {
	_, err := s.pool.Exec(context.Background(), `DELETE FROM comments WHERE id = $1`, id)
	return err
}

// Attachments
func (s *Storage) CreateAttachment(a *models.Attachment) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO attachments (id, card_id, user_id, file_name, file_path, file_size, mime_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.CardID, a.UserID, a.FileName, a.FilePath, a.FileSize, a.MimeType, a.CreatedAt,
	)
	return err
}

func (s *Storage) GetAttachmentsByCard(cardID string) ([]models.Attachment, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, card_id, user_id, file_name, file_size, mime_type, created_at
		FROM attachments WHERE card_id = $1 ORDER BY created_at ASC`,
		cardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var atts []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.CardID, &a.UserID, &a.FileName, &a.FileSize, &a.MimeType, &a.CreatedAt); err != nil {
			return nil, err
		}
		atts = append(atts, a)
	}
	if atts == nil {
		atts = []models.Attachment{}
	}
	return atts, nil
}

func (s *Storage) GetAttachment(id string) (*models.Attachment, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, card_id, user_id, file_name, file_path, file_size, mime_type, created_at
		FROM attachments WHERE id = $1`, id,
	)
	var a models.Attachment
	err := row.Scan(&a.ID, &a.CardID, &a.UserID, &a.FileName, &a.FilePath, &a.FileSize, &a.MimeType, &a.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (s *Storage) DeleteAttachment(id string) error {
	_, err := s.pool.Exec(context.Background(), `DELETE FROM attachments WHERE id = $1`, id)
	return err
}

// Users
func (s *Storage) CreateUser(user *models.User) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (s *Storage) GetUserByEmail(email string) (*models.User, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`, email,
	)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetUserByID(id string) (*models.User, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`, id,
	)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetBoardOwner(boardID string) (string, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT user_id FROM boards WHERE id = $1`, boardID,
	)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return userID, nil
}

func (s *Storage) GetListBoardID(listID string) (string, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT board_id FROM lists WHERE id = $1`, listID,
	)
	var boardID string
	err := row.Scan(&boardID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return boardID, nil
}

func (s *Storage) GetCardListID(cardID string) (string, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT list_id FROM cards WHERE id = $1`, cardID,
	)
	var listID string
	err := row.Scan(&listID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return listID, nil
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Storage) GetUserInfo(id string) (*UserInfo, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, username, email FROM users WHERE id = $1`, id,
	)
	var u UserInfo
	err := row.Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
