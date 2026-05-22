package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"artist-distribution-study/backend/internal/domain"
	"artist-distribution-study/backend/internal/repository"
)

// This file models release submission and store delivery status changes.

type ReleaseService struct {
	artists  repository.ArtistRepository
	tracks   repository.TrackRepository
	stores   repository.StoreRepository
	releases repository.ReleaseRepository
}

type CreateReleaseRequestInput struct {
	ArtistID    string   `json:"artistId"`
	TrackID     string   `json:"trackId"`
	Title       string   `json:"title"`
	ReleaseDate string   `json:"releaseDate"`
	StoreIDs    []string `json:"storeIds"`
}

type UpdateReleaseStatusInput struct {
	ReleaseID   string
	Status      *domain.ReleaseStatus
	StoreID     string
	StoreStatus *domain.StoreDeliveryStatus
}

func NewReleaseService(artists repository.ArtistRepository, tracks repository.TrackRepository, stores repository.StoreRepository, releases repository.ReleaseRepository) *ReleaseService {
	return &ReleaseService{
		artists:  artists,
		tracks:   tracks,
		stores:   stores,
		releases: releases,
	}
}

func (s *ReleaseService) ListStores(ctx context.Context) ([]domain.Store, error) {
	return s.stores.ListStores(ctx)
}

func (s *ReleaseService) CreateReleaseRequest(ctx context.Context, input CreateReleaseRequestInput) (*domain.ReleaseRequest, error) {
	if _, err := s.artists.GetArtist(ctx, input.ArtistID); err != nil {
		return nil, err
	}

	track, err := s.tracks.GetTrack(ctx, input.TrackID)
	if err != nil {
		return nil, err
	}

	if track.ArtistID != input.ArtistID {
		return nil, errors.New("track does not belong to artist")
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = track.Title + " - Single"
	}

	releaseDate, err := time.Parse("2006-01-02", input.ReleaseDate)
	if err != nil {
		return nil, errors.New("releaseDate must use YYYY-MM-DD")
	}

	if len(input.StoreIDs) == 0 {
		return nil, errors.New("at least one store is required")
	}

	stores, err := s.stores.ListStores(ctx)
	if err != nil {
		return nil, err
	}

	storeByID := make(map[string]domain.Store, len(stores))
	for _, store := range stores {
		storeByID[store.ID] = store
	}

	seenStores := make(map[string]bool, len(input.StoreIDs))
	deliveries := make([]domain.StoreDelivery, 0, len(input.StoreIDs))
	for _, storeID := range input.StoreIDs {
		if seenStores[storeID] {
			continue
		}

		store, ok := storeByID[storeID]
		if !ok {
			return nil, fmt.Errorf("store %s: %w", storeID, repository.ErrNotFound)
		}

		deliveries = append(deliveries, domain.StoreDelivery{
			StoreID:   store.ID,
			StoreName: store.Name,
			Status:    domain.StoreDeliveryStatusPending,
		})
		seenStores[storeID] = true
	}

	request := domain.ReleaseRequest{
		ArtistID:        input.ArtistID,
		TrackID:         input.TrackID,
		Title:           title,
		ReleaseDate:     releaseDate,
		Status:          domain.ReleaseStatusSubmitted,
		StoreDeliveries: deliveries,
	}

	return s.releases.CreateReleaseRequest(ctx, request)
}

func (s *ReleaseService) ListReleaseRequests(ctx context.Context, artistID string) ([]domain.ReleaseRequest, error) {
	if _, err := s.artists.GetArtist(ctx, artistID); err != nil {
		return nil, err
	}

	return s.releases.ListReleaseRequests(ctx, artistID)
}

func (s *ReleaseService) UpdateReleaseStatus(ctx context.Context, input UpdateReleaseStatusInput) (*domain.ReleaseRequest, error) {
	if input.Status == nil && input.StoreStatus == nil {
		return nil, errors.New("status or storeStatus is required")
	}

	if input.Status != nil && !isValidReleaseStatus(*input.Status) {
		return nil, errors.New("invalid release status")
	}

	if input.StoreStatus != nil {
		if strings.TrimSpace(input.StoreID) == "" {
			return nil, errors.New("storeId is required when updating storeStatus")
		}

		if !isValidStoreDeliveryStatus(*input.StoreStatus) {
			return nil, errors.New("invalid store delivery status")
		}
	}

	return s.releases.UpdateReleaseStatus(ctx, input.ReleaseID, input.Status, input.StoreID, input.StoreStatus)
}

func isValidReleaseStatus(status domain.ReleaseStatus) bool {
	switch status {
	case domain.ReleaseStatusDraft,
		domain.ReleaseStatusSubmitted,
		domain.ReleaseStatusInReview,
		domain.ReleaseStatusApproved,
		domain.ReleaseStatusDelivered,
		domain.ReleaseStatusLive,
		domain.ReleaseStatusRejected:
		return true
	default:
		return false
	}
}

func isValidStoreDeliveryStatus(status domain.StoreDeliveryStatus) bool {
	switch status {
	case domain.StoreDeliveryStatusPending,
		domain.StoreDeliveryStatusProcessing,
		domain.StoreDeliveryStatusDelivered,
		domain.StoreDeliveryStatusLive,
		domain.StoreDeliveryStatusFailed:
		return true
	default:
		return false
	}
}
