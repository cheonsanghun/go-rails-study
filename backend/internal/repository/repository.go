package repository

import (
	"context"
	"errors"

	"artist-distribution-study/backend/internal/domain"
)

// This file declares repository interfaces so PostgreSQL can replace memory storage later.

var ErrNotFound = errors.New("not found")

type ArtistRepository interface {
	ListArtists(ctx context.Context) ([]domain.Artist, error)
	GetArtist(ctx context.Context, artistID string) (*domain.Artist, error)
}

type TrackRepository interface {
	CreateTrack(ctx context.Context, track domain.Track) (*domain.Track, error)
	ListTracks(ctx context.Context, artistID string) ([]domain.Track, error)
	GetTrack(ctx context.Context, trackID string) (*domain.Track, error)
}

type StoreRepository interface {
	ListStores(ctx context.Context) ([]domain.Store, error)
}

type ReleaseRepository interface {
	CreateReleaseRequest(ctx context.Context, request domain.ReleaseRequest) (*domain.ReleaseRequest, error)
	ListReleaseRequests(ctx context.Context, artistID string) ([]domain.ReleaseRequest, error)
	GetReleaseRequest(ctx context.Context, releaseID string) (*domain.ReleaseRequest, error)
	UpdateReleaseStatus(ctx context.Context, releaseID string, status *domain.ReleaseStatus, storeID string, storeStatus *domain.StoreDeliveryStatus) (*domain.ReleaseRequest, error)
}

type SalesRepository interface {
	ListSalesReports(ctx context.Context, artistID string) ([]domain.SalesReport, error)
}
