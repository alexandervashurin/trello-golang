package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexandervashurin/trello-golang/handlers"
	"github.com/alexandervashurin/trello-golang/storage"
)

func main() {
	db, err := storage.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := storage.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	store := storage.NewStorage(db)
	handler := handlers.NewHandler(store)
	authHandler := handlers.NewAuthHandler(store)
	commentHandler := handlers.NewCommentHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", handlers.LoginLimiter(authHandler.Register))
	mux.HandleFunc("POST /api/login", handlers.LoginLimiter(authHandler.Login))
	mux.HandleFunc("GET /api/me", handlers.AuthMiddleware(authHandler.Me))

	mux.HandleFunc("POST /api/boards", handlers.AuthMiddleware(handler.CreateBoard))
	mux.HandleFunc("GET /api/boards", handlers.AuthMiddleware(handler.GetAllBoards))
	mux.HandleFunc("GET /api/boards/public", handler.GetPublicBoards)
	mux.HandleFunc("GET /api/board", handlers.OptionalAuth(handler.GetBoard))

	mux.HandleFunc("POST /api/lists", handlers.AuthMiddleware(handler.CreateList))
	mux.HandleFunc("GET /api/lists", handlers.OptionalAuth(handler.GetListsByBoard))

	mux.HandleFunc("POST /api/cards", handlers.AuthMiddleware(handler.CreateCard))
	mux.HandleFunc("GET /api/cards", handlers.OptionalAuth(handler.GetCardsByList))
	mux.HandleFunc("PATCH /api/card", handlers.AuthMiddleware(handler.MoveCard))
	mux.HandleFunc("DELETE /api/card", handlers.AuthMiddleware(handler.DeleteCard))

	mux.HandleFunc("POST /api/comments", handlers.AuthMiddleware(commentHandler.CreateComment))
	mux.HandleFunc("GET /api/comments", handlers.OptionalAuth(commentHandler.GetComments))
	mux.HandleFunc("DELETE /api/comment", handlers.AuthMiddleware(commentHandler.DeleteComment))

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
