package domain

import "time"

// This file defines the core music distribution domain models used by all layers.

type Artist struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	LabelName    string    `json:"labelName"`
	Country      string    `json:"country"`
	PrimaryGenre string    `json:"primaryGenre"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Track struct {
	ID              string    `json:"id"`
	ArtistID        string    `json:"artistId"`
	Title           string    `json:"title"`
	ISRC            string    `json:"isrc"`
	DurationSeconds int       `json:"durationSeconds"`
	Genre           string    `json:"genre"`
	Language        string    `json:"language"`
	Explicit        bool      `json:"explicit"`
	AudioFileName   string    `json:"audioFileName"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Store struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Territory        string `json:"territory"`
	SupportsPreorder bool   `json:"supportsPreorder"`
}

type ReleaseStatus string

const (
	ReleaseStatusDraft     ReleaseStatus = "draft"
	ReleaseStatusSubmitted ReleaseStatus = "submitted"
	ReleaseStatusInReview  ReleaseStatus = "in_review"
	ReleaseStatusApproved  ReleaseStatus = "approved"
	ReleaseStatusDelivered ReleaseStatus = "delivered"
	ReleaseStatusLive      ReleaseStatus = "live"
	ReleaseStatusRejected  ReleaseStatus = "rejected"
)

type StoreDeliveryStatus string

const (
	StoreDeliveryStatusPending    StoreDeliveryStatus = "pending"
	StoreDeliveryStatusProcessing StoreDeliveryStatus = "processing"
	StoreDeliveryStatusDelivered  StoreDeliveryStatus = "delivered"
	StoreDeliveryStatusLive       StoreDeliveryStatus = "live"
	StoreDeliveryStatusFailed     StoreDeliveryStatus = "failed"
)

type StoreDelivery struct {
	StoreID       string              `json:"storeId"`
	StoreName     string              `json:"storeName"`
	Status        StoreDeliveryStatus `json:"status"`
	LastUpdatedAt time.Time           `json:"lastUpdatedAt"`
}

type ReleaseRequest struct {
	ID              string          `json:"id"`
	ArtistID        string          `json:"artistId"`
	TrackID         string          `json:"trackId"`
	Title           string          `json:"title"`
	ReleaseDate     time.Time       `json:"releaseDate"`
	Status          ReleaseStatus   `json:"status"`
	StoreDeliveries []StoreDelivery `json:"storeDeliveries"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type SalesReport struct {
	ID              string    `json:"id"`
	ArtistID        string    `json:"artistId"`
	TrackID         string    `json:"trackId"`
	TrackTitle      string    `json:"trackTitle"`
	StoreID         string    `json:"storeId"`
	StoreName       string    `json:"storeName"`
	Period          string    `json:"period"`
	Country         string    `json:"country"`
	Currency        string    `json:"currency"`
	Units           int       `json:"units"`
	NetRevenueMinor int       `json:"netRevenueMinor"`
	ReportedAt      time.Time `json:"reportedAt"`
}
