package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexandervashurin/trello-golang/models"
	"github.com/alexandervashurin/trello-golang/storage"
	"github.com/alexandervashurin/trello-golang/utils"
	"github.com/google/uuid"
)

const maxUploadSize = 10 << 20

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true,
	".pdf": true,
	".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
	".txt": true, ".csv": true,
	".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
	".mp3": true, ".wav": true, ".ogg": true,
	".mp4": true, ".avi": true, ".mov": true, ".mkv": true,
}

var allowedMimePrefixes = []string{
	"image/",
	"text/",
	"application/pdf",
	"application/msword",
	"application/vnd.openxmlformats-officedocument",
	"application/vnd.ms-",
	"application/zip", "application/x-zip", "application/x-rar", "application/x-7z",
	"application/gzip", "application/x-tar",
	"audio/",
	"video/",
}

type AttachmentHandler struct {
	storage   *storage.Storage
	uploadDir string
}

func NewAttachmentHandler(s *storage.Storage, uploadDir string) *AttachmentHandler {
	os.MkdirAll(uploadDir, 0755)
	return &AttachmentHandler{storage: s, uploadDir: uploadDir}
}

func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "File too large (max 10MB)")
		return
	}

	cardID := r.FormValue("card_id")
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

	if !board.IsPublic && board.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if !allowedExtensions[ext] {
		utils.RespondWithError(w, http.StatusBadRequest, "File type not allowed")
		return
	}

	fileID := uuid.New().String()
	savedName := fileID + ext
	savedPath := filepath.Join(h.uploadDir, savedName)

	dst, err := os.Create(savedPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(savedPath)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to write file")
		return
	}

	mime := header.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}

	allowed := false
	for _, prefix := range allowedMimePrefixes {
		if len(mime) >= len(prefix) && mime[:len(prefix)] == prefix {
			allowed = true
			break
		}
	}
	if !allowed {
		os.Remove(savedPath)
		utils.RespondWithError(w, http.StatusBadRequest, "File type not allowed")
		return
	}

	now := time.Now()
	att := &models.Attachment{
		ID:        fileID,
		CardID:    cardID,
		UserID:    userID,
		FileName:  header.Filename,
		FilePath:  savedPath,
		FileSize:  written,
		MimeType:  mime,
		CreatedAt: now,
	}

	if err := h.storage.CreateAttachment(att); err != nil {
		os.Remove(savedPath)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save metadata")
		return
	}

	att.FilePath = ""
	utils.RespondWithSuccess(w, att)
}

func (h *AttachmentHandler) GetAttachments(w http.ResponseWriter, r *http.Request) {
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

	atts, err := h.storage.GetAttachmentsByCard(cardID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondWithSuccess(w, atts)
}

func (h *AttachmentHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	att, err := h.storage.GetAttachment(id)
	if err != nil || att == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Attachment not found")
		return
	}

	boardID, err := h.storage.GetCardBoardID(att.CardID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	board, err := h.storage.GetBoard(boardID)
	if err != nil || board == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Board not found")
		return
	}

	if board.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	os.Remove(att.FilePath)

	if err := h.storage.DeleteAttachment(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete attachment")
		return
	}

	utils.RespondWithSuccess(w, map[string]string{"message": "Attachment deleted"})
}

func (h *AttachmentHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	att, err := h.storage.GetAttachment(id)
	if err != nil || att == nil {
		http.NotFound(w, r)
		return
	}

	boardID, err := h.storage.GetCardBoardID(att.CardID)
	if err != nil || boardID == "" {
		http.NotFound(w, r)
		return
	}

	board, err := h.storage.GetBoard(boardID)
	if err != nil || board == nil {
		http.NotFound(w, r)
		return
	}

	userID, _ := r.Context().Value(UserIDKey).(string)
	if !board.IsPublic && board.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, att.FileName))
	w.Header().Set("Content-Type", att.MimeType)
	http.ServeFile(w, r, att.FilePath)
}
