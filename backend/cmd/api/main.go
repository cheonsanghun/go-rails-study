package main

import (
	"log"
	"net/http"

	"artist-distribution-study/backend/internal/handler"
	"artist-distribution-study/backend/internal/repository"
	"artist-distribution-study/backend/internal/service"
)

// This file starts the local REST API with an in-memory repository.
func main() {
	repo := repository.NewMemoryRepository()

	artistService := service.NewArtistService(repo)
	trackService := service.NewTrackService(repo, repo)
	releaseService := service.NewReleaseService(repo, repo, repo, repo)
	salesService := service.NewSalesService(repo, repo)

	router := handler.NewRouter(artistService, trackService, releaseService, salesService)

	addr := ":8080"
	log.Printf("artist distribution API listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
