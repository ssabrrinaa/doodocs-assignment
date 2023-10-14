// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"test/internal/handlers"
	"test/internal/repository"
	"test/internal/service"
)

const (
	Addr = ":4000"
)

func main() {
	r := repository.NewArchiveRepository()
	s := service.NewArchiveService(r)
	h := handlers.NewArchiveHandler(s)
	// archiveHandler := handlers.NewArchiveHandler()

	// Регистрируем хендлер
	http.HandleFunc("/api/archive/information", h.HandleArchiveInformation)

	// http.ListenAndServe(":8080", nil)

	fmt.Println("Server starting on http://localhost:4000")
	if err := http.ListenAndServe(Addr, nil); err != nil {
		log.Fatal(err)
	}
}
