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

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/api/archive/information", h.ArchiveInformationHandler)
	http.HandleFunc("/api/archive/files", h.CreateArchiveHandler)

	fmt.Println("Server starting on http://localhost:4000")
	if err := http.ListenAndServe(Addr, nil); err != nil {
		log.Fatal(err)
	}
}
