package service

import (
	"context"

	"artist-distribution-study/backend/internal/domain"
	"artist-distribution-study/backend/internal/repository"
)

// This file contains artist use cases and keeps handlers independent from storage.

type ArtistService struct {
	artists repository.ArtistRepository
}

func NewArtistService(artists repository.ArtistRepository) *ArtistService {
	return &ArtistService{artists: artists}
}

func (s *ArtistService) ListArtists(ctx context.Context) ([]domain.Artist, error) {
	return s.artists.ListArtists(ctx)
}

func (s *ArtistService) GetArtist(ctx context.Context, artistID string) (*domain.Artist, error) {
	return s.artists.GetArtist(ctx, artistID)
}
