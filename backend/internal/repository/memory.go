package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"artist-distribution-study/backend/internal/domain"
)

// This file provides a beginner-friendly in-memory repository for local development.

type MemoryRepository struct {
	mu sync.RWMutex

	artists         map[string]domain.Artist
	tracks          map[string]domain.Track
	stores          map[string]domain.Store
	releaseRequests map[string]domain.ReleaseRequest
	salesReports    []domain.SalesReport

	nextTrackNumber   int
	nextReleaseNumber int
}

func NewMemoryRepository() *MemoryRepository {
	now := time.Now().UTC().Truncate(time.Second)

	repo := &MemoryRepository{
		artists: map[string]domain.Artist{
			"art_1001": {
				ID:           "art_1001",
				Name:         "Hikari Lane",
				LabelName:    "Sunrise Room",
				Country:      "JP",
				PrimaryGenre: "Electronic",
				CreatedAt:    now.AddDate(0, -4, 0),
			},
			"art_1002": {
				ID:           "art_1002",
				Name:         "North Harbor Kids",
				LabelName:    "Self-Released",
				Country:      "JP",
				PrimaryGenre: "Indie Pop",
				CreatedAt:    now.AddDate(0, -2, -10),
			},
			"art_1003": {
				ID:           "art_1003",
				Name:         "Mina Soto",
				LabelName:    "Paper Crane Records",
				Country:      "JP",
				PrimaryGenre: "R&B",
				CreatedAt:    now.AddDate(0, -1, -3),
			},
		},
		tracks: map[string]domain.Track{
			"trk_1001": {
				ID:              "trk_1001",
				ArtistID:        "art_1001",
				Title:           "Late Night Conveyor",
				ISRC:            "JPA1B2600010",
				DurationSeconds: 214,
				Genre:           "Electronic",
				Language:        "ja",
				Explicit:        false,
				AudioFileName:   "late-night-conveyor.wav",
				CreatedAt:       now.AddDate(0, -1, -12),
			},
			"trk_1002": {
				ID:              "trk_1002",
				ArtistID:        "art_1002",
				Title:           "Harbor Light Demo",
				ISRC:            "JPA1B2600011",
				DurationSeconds: 188,
				Genre:           "Indie Pop",
				Language:        "ja",
				Explicit:        false,
				AudioFileName:   "harbor-light-demo.wav",
				CreatedAt:       now.AddDate(0, -1, -4),
			},
		},
		stores: map[string]domain.Store{
			"spotify": {
				ID:               "spotify",
				Name:             "Spotify",
				Territory:        "Global",
				SupportsPreorder: false,
			},
			"apple_music": {
				ID:               "apple_music",
				Name:             "Apple Music",
				Territory:        "Global",
				SupportsPreorder: true,
			},
			"line_music": {
				ID:               "line_music",
				Name:             "LINE MUSIC",
				Territory:        "Japan",
				SupportsPreorder: false,
			},
			"youtube_music": {
				ID:               "youtube_music",
				Name:             "YouTube Music",
				Territory:        "Global",
				SupportsPreorder: false,
			},
			"amazon_music": {
				ID:               "amazon_music",
				Name:             "Amazon Music",
				Territory:        "Global",
				SupportsPreorder: true,
			},
		},
		releaseRequests: map[string]domain.ReleaseRequest{
			"rel_1001": {
				ID:          "rel_1001",
				ArtistID:    "art_1001",
				TrackID:     "trk_1001",
				Title:       "Late Night Conveyor - Single",
				ReleaseDate: now.AddDate(0, 0, 21),
				Status:      domain.ReleaseStatusInReview,
				StoreDeliveries: []domain.StoreDelivery{
					{StoreID: "spotify", StoreName: "Spotify", Status: domain.StoreDeliveryStatusProcessing, LastUpdatedAt: now.AddDate(0, 0, -2)},
					{StoreID: "apple_music", StoreName: "Apple Music", Status: domain.StoreDeliveryStatusDelivered, LastUpdatedAt: now.AddDate(0, 0, -1)},
					{StoreID: "line_music", StoreName: "LINE MUSIC", Status: domain.StoreDeliveryStatusPending, LastUpdatedAt: now.AddDate(0, 0, -2)},
				},
				CreatedAt: now.AddDate(0, 0, -7),
				UpdatedAt: now.AddDate(0, 0, -1),
			},
		},
		salesReports: []domain.SalesReport{
			{
				ID:              "sale_1001",
				ArtistID:        "art_1001",
				TrackID:         "trk_1001",
				TrackTitle:      "Late Night Conveyor",
				StoreID:         "spotify",
				StoreName:       "Spotify",
				Period:          "2026-04",
				Country:         "JP",
				Currency:        "JPY",
				Units:           2840,
				NetRevenueMinor: 18420,
				ReportedAt:      now.AddDate(0, 0, -15),
			},
			{
				ID:              "sale_1002",
				ArtistID:        "art_1001",
				TrackID:         "trk_1001",
				TrackTitle:      "Late Night Conveyor",
				StoreID:         "apple_music",
				StoreName:       "Apple Music",
				Period:          "2026-04",
				Country:         "JP",
				Currency:        "JPY",
				Units:           1290,
				NetRevenueMinor: 12600,
				ReportedAt:      now.AddDate(0, 0, -15),
			},
			{
				ID:              "sale_1003",
				ArtistID:        "art_1002",
				TrackID:         "trk_1002",
				TrackTitle:      "Harbor Light Demo",
				StoreID:         "line_music",
				StoreName:       "LINE MUSIC",
				Period:          "2026-04",
				Country:         "JP",
				Currency:        "JPY",
				Units:           740,
				NetRevenueMinor: 5120,
				ReportedAt:      now.AddDate(0, 0, -15),
			},
		},
		nextTrackNumber:   1003,
		nextReleaseNumber: 1002,
	}

	return repo
}

func (r *MemoryRepository) ListArtists(ctx context.Context) ([]domain.Artist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artists := make([]domain.Artist, 0, len(r.artists))
	for _, artist := range r.artists {
		artists = append(artists, artist)
	}

	sort.Slice(artists, func(i, j int) bool {
		return artists[i].Name < artists[j].Name
	})

	return artists, nil
}

func (r *MemoryRepository) GetArtist(ctx context.Context, artistID string) (*domain.Artist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artist, ok := r.artists[artistID]
	if !ok {
		return nil, ErrNotFound
	}

	return &artist, nil
}

func (r *MemoryRepository) CreateTrack(ctx context.Context, track domain.Track) (*domain.Track, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.artists[track.ArtistID]; !ok {
		return nil, ErrNotFound
	}

	now := time.Now().UTC().Truncate(time.Second)
	track.ID = fmt.Sprintf("trk_%04d", r.nextTrackNumber)
	track.CreatedAt = now
	r.nextTrackNumber++

	r.tracks[track.ID] = track

	return &track, nil
}

func (r *MemoryRepository) ListTracks(ctx context.Context, artistID string) ([]domain.Track, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tracks := make([]domain.Track, 0)
	for _, track := range r.tracks {
		if artistID == "" || track.ArtistID == artistID {
			tracks = append(tracks, track)
		}
	}

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].CreatedAt.After(tracks[j].CreatedAt)
	})

	return tracks, nil
}

func (r *MemoryRepository) GetTrack(ctx context.Context, trackID string) (*domain.Track, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	track, ok := r.tracks[trackID]
	if !ok {
		return nil, ErrNotFound
	}

	return &track, nil
}

func (r *MemoryRepository) ListStores(ctx context.Context) ([]domain.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stores := make([]domain.Store, 0, len(r.stores))
	for _, store := range r.stores {
		stores = append(stores, store)
	}

	sort.Slice(stores, func(i, j int) bool {
		return stores[i].Name < stores[j].Name
	})

	return stores, nil
}

func (r *MemoryRepository) CreateReleaseRequest(ctx context.Context, request domain.ReleaseRequest) (*domain.ReleaseRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.artists[request.ArtistID]; !ok {
		return nil, ErrNotFound
	}

	track, ok := r.tracks[request.TrackID]
	if !ok || track.ArtistID != request.ArtistID {
		return nil, ErrNotFound
	}

	now := time.Now().UTC().Truncate(time.Second)
	request.ID = fmt.Sprintf("rel_%04d", r.nextReleaseNumber)
	request.Status = domain.ReleaseStatusSubmitted
	request.CreatedAt = now
	request.UpdatedAt = now

	for i := range request.StoreDeliveries {
		request.StoreDeliveries[i].LastUpdatedAt = now
	}

	r.nextReleaseNumber++
	r.releaseRequests[request.ID] = request

	return &request, nil
}

func (r *MemoryRepository) ListReleaseRequests(ctx context.Context, artistID string) ([]domain.ReleaseRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	requests := make([]domain.ReleaseRequest, 0)
	for _, request := range r.releaseRequests {
		if artistID == "" || request.ArtistID == artistID {
			requests = append(requests, request)
		}
	}

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].CreatedAt.After(requests[j].CreatedAt)
	})

	return requests, nil
}

func (r *MemoryRepository) GetReleaseRequest(ctx context.Context, releaseID string) (*domain.ReleaseRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	request, ok := r.releaseRequests[releaseID]
	if !ok {
		return nil, ErrNotFound
	}

	return &request, nil
}

func (r *MemoryRepository) UpdateReleaseStatus(ctx context.Context, releaseID string, status *domain.ReleaseStatus, storeID string, storeStatus *domain.StoreDeliveryStatus) (*domain.ReleaseRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	request, ok := r.releaseRequests[releaseID]
	if !ok {
		return nil, ErrNotFound
	}

	if status != nil {
		request.Status = *status
	}

	if storeStatus != nil {
		updated := false
		for i := range request.StoreDeliveries {
			if request.StoreDeliveries[i].StoreID == storeID {
				request.StoreDeliveries[i].Status = *storeStatus
				request.StoreDeliveries[i].LastUpdatedAt = time.Now().UTC().Truncate(time.Second)
				updated = true
				break
			}
		}

		if !updated {
			return nil, ErrNotFound
		}
	}

	request.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	r.releaseRequests[releaseID] = request

	return &request, nil
}

func (r *MemoryRepository) ListSalesReports(ctx context.Context, artistID string) ([]domain.SalesReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reports := make([]domain.SalesReport, 0)
	for _, report := range r.salesReports {
		if artistID == "" || report.ArtistID == artistID {
			reports = append(reports, report)
		}
	}

	sort.Slice(reports, func(i, j int) bool {
		if reports[i].Period == reports[j].Period {
			return reports[i].StoreName < reports[j].StoreName
		}
		return reports[i].Period > reports[j].Period
	})

	return reports, nil
}
