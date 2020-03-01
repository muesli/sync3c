package main

import (
	"encoding/json"
	"net/http"
)

// Recording represents a single recording
type Recording struct {
	ConferenceURL string `json:"conference_url"`
	EventURL      string `json:"event_url"`
	Filename      string `json:"filename"`
	Folder        string `json:"folder"`
	Height        int64  `json:"height"`
	HighQuality   bool   `json:"high_quality"`
	Language      string `json:"language"`
	Length        int64  `json:"length"`
	MimeType      string `json:"mime_type"`
	RecordingURL  string `json:"recording_url"`
	Size          int64  `json:"size"`
	State         string `json:"state"`
	UpdatedAt     string `json:"updated_at"`
	URL           string `json:"url"`
	Width         int64  `json:"width"`
}

// Media represents a single talk
type Media struct {
	ConferenceURL    string      `json:"conference_url"`
	Date             string      `json:"date"`
	Description      string      `json:"description"`
	FrontendLink     string      `json:"frontend_link"`
	GUID             string      `json:"guid"`
	Length           int64       `json:"length"`
	Link             string      `json:"link"`
	OriginalLanguage string      `json:"original_language"`
	Persons          []string    `json:"persons"`
	PosterURL        string      `json:"poster_url"`
	Recordings       []Recording `json:"recordings"`
	ReleaseDate      string      `json:"release_date"`
	Slug             string      `json:"slug"`
	Subtitle         string      `json:"subtitle"`
	Tags             []string    `json:"tags"`
	ThumbURL         string      `json:"thumb_url"`
	Title            string      `json:"title"`
	UpdatedAt        string      `json:"updated_at"`
	URL              string      `json:"url"`
}

func findMedia(url string) (Media, error) {
	media := Media{}

	r, err := http.Get(url)
	if err != nil {
		return media, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&media)
	return media, err
}
