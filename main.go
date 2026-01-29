package main

import (
	"log"
	"net/http"

	"github.com/alexandervashurin/trello-golang/handlers"
	"github.com/alexandervashurin/trello-golang/storage"
)

func main() {
	// Инициализация хранилища
	store := storage.NewStorage()
	handler := handlers.NewHandler(store)

	// Настройка маршрутов
	mux := http.NewServeMux()

	// Board routes
	mux.HandleFunc("POST /api/boards", handler.CreateBoard)
	mux.HandleFunc("GET /api/boards", handler.GetAllBoards)
	mux.HandleFunc("GET /api/board", handler.GetBoard)
	mux.HandleFunc("DELETE /api/board", handler.DeleteBoard)

	// List routes
	mux.HandleFunc("POST /api/lists", handler.CreateList)
	mux.HandleFunc("GET /api/lists", handler.GetListsByBoard)
	mux.HandleFunc("DELETE /api/list", handler.DeleteList)

	// Card routes
	mux.HandleFunc("POST /api/cards", handler.CreateCard)
	mux.HandleFunc("GET /api/cards", handler.GetCardsByList)
	mux.HandleFunc("DELETE /api/card", handler.DeleteCard)

	// Запуск сервера
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
