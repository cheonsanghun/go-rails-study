package service

import (
	"context"

	"artist-distribution-study/backend/internal/domain"
	"artist-distribution-study/backend/internal/repository"
)

// This file exposes read-only revenue report use cases for the frontend dashboard.

type SalesService struct {
	artists repository.ArtistRepository
	sales   repository.SalesRepository
}

func NewSalesService(artists repository.ArtistRepository, sales repository.SalesRepository) *SalesService {
	return &SalesService{artists: artists, sales: sales}
}

func (s *SalesService) ListSalesReports(ctx context.Context, artistID string) ([]domain.SalesReport, error) {
	if _, err := s.artists.GetArtist(ctx, artistID); err != nil {
		return nil, err
	}

	return s.sales.ListSalesReports(ctx, artistID)
}
