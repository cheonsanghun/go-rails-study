package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"artist-distribution-study/backend/internal/domain"
	"artist-distribution-study/backend/internal/repository"
	"artist-distribution-study/backend/internal/service"
)

// This file maps REST endpoints to service methods and owns HTTP concerns.

type Handler struct {
	artists  *service.ArtistService
	tracks   *service.TrackService
	releases *service.ReleaseService
	sales    *service.SalesService
}

type updateReleaseStatusPayload struct {
	Status      *domain.ReleaseStatus       `json:"status"`
	StoreID     string                      `json:"storeId"`
	StoreStatus *domain.StoreDeliveryStatus `json:"storeStatus"`
}

func NewRouter(artists *service.ArtistService, tracks *service.TrackService, releases *service.ReleaseService, sales *service.SalesService) http.Handler {
	h := &Handler{
		artists:  artists,
		tracks:   tracks,
		releases: releases,
		sales:    sales,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/api/stores", h.handleStores)
	mux.HandleFunc("/api/artists", h.handleArtists)
	mux.HandleFunc("/api/artists/", h.handleArtistSubroutes)
	mux.HandleFunc("/api/release-requests/", h.handleReleaseSubroutes)

	return withCORS(mux)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleStores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	stores, err := h.releases.ListStores(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"stores": stores})
}

func (h *Handler) handleArtists(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/artists" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	artists, err := h.artists.ListArtists(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"artists": artists})
}

func (h *Handler) handleArtistSubroutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/artists/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	artistID := parts[0]
	if len(parts) == 1 {
		h.handleArtistDetail(w, r, artistID)
		return
	}

	switch parts[1] {
	case "tracks":
		h.handleTracks(w, r, artistID)
	case "release-requests":
		h.handleReleaseRequests(w, r, artistID)
	case "sales-reports":
		h.handleSalesReports(w, r, artistID)
	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

func (h *Handler) handleArtistDetail(w http.ResponseWriter, r *http.Request, artistID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	artist, err := h.artists.GetArtist(r.Context(), artistID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"artist": artist})
}

func (h *Handler) handleTracks(w http.ResponseWriter, r *http.Request, artistID string) {
	switch r.Method {
	case http.MethodGet:
		tracks, err := h.tracks.ListTracks(r.Context(), artistID)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"tracks": tracks})
	case http.MethodPost:
		var input service.CreateTrackInput
		if err := readJSON(r, &input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		input.ArtistID = artistID
		track, err := h.tracks.CreateTrack(r.Context(), input)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{"track": track})
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleReleaseRequests(w http.ResponseWriter, r *http.Request, artistID string) {
	switch r.Method {
	case http.MethodGet:
		requests, err := h.releases.ListReleaseRequests(r.Context(), artistID)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"releaseRequests": requests})
	case http.MethodPost:
		var input service.CreateReleaseRequestInput
		if err := readJSON(r, &input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		input.ArtistID = artistID
		request, err := h.releases.CreateReleaseRequest(r.Context(), input)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{"releaseRequest": request})
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleSalesReports(w http.ResponseWriter, r *http.Request, artistID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	reports, err := h.sales.ListSalesReports(r.Context(), artistID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"salesReports": reports})
}

func (h *Handler) handleReleaseSubroutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/release-requests/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] != "status" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	if r.Method != http.MethodPatch {
		writeMethodNotAllowed(w)
		return
	}

	var payload updateReleaseStatusPayload
	if err := readJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	request, err := h.releases.UpdateReleaseStatus(r.Context(), service.UpdateReleaseStatusInput{
		ReleaseID:   parts[0],
		Status:      payload.Status,
		StoreID:     payload.StoreID,
		StoreStatus: payload.StoreStatus,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"releaseRequest": request})
}

func readJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
