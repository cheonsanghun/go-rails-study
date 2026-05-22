package service

import (
	"context"
	"errors"
	"strings"

	"artist-distribution-study/backend/internal/domain"
	"artist-distribution-study/backend/internal/repository"
)

// This file contains track registration rules before data is saved.

type TrackService struct {
	artists repository.ArtistRepository
	tracks  repository.TrackRepository
}

type CreateTrackInput struct {
	ArtistID        string `json:"artistId"`
	Title           string `json:"title"`
	ISRC            string `json:"isrc"`
	DurationSeconds int    `json:"durationSeconds"`
	Genre           string `json:"genre"`
	Language        string `json:"language"`
	Explicit        bool   `json:"explicit"`
	AudioFileName   string `json:"audioFileName"`
}

func NewTrackService(artists repository.ArtistRepository, tracks repository.TrackRepository) *TrackService {
	return &TrackService{artists: artists, tracks: tracks}
}

func (s *TrackService) CreateTrack(ctx context.Context, input CreateTrackInput) (*domain.Track, error) {
	artist, err := s.artists.GetArtist(ctx, input.ArtistID)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, errors.New("title is required")
	}

	if input.DurationSeconds <= 0 {
		return nil, errors.New("durationSeconds must be greater than zero")
	}

	genre := strings.TrimSpace(input.Genre)
	if genre == "" {
		genre = artist.PrimaryGenre
	}

	language := strings.TrimSpace(input.Language)
	if language == "" {
		language = "ja"
	}

	track := domain.Track{
		ArtistID:        input.ArtistID,
		Title:           title,
		ISRC:            strings.TrimSpace(input.ISRC),
		DurationSeconds: input.DurationSeconds,
		Genre:           genre,
		Language:        language,
		Explicit:        input.Explicit,
		AudioFileName:   strings.TrimSpace(input.AudioFileName),
	}

	return s.tracks.CreateTrack(ctx, track)
}

func (s *TrackService) ListTracks(ctx context.Context, artistID string) ([]domain.Track, error) {
	if _, err := s.artists.GetArtist(ctx, artistID); err != nil {
		return nil, err
	}

	return s.tracks.ListTracks(ctx, artistID)
}
