package main

import (
	"encoding/json"
	"net/http"
)

type Event struct {
	Acronym     string `json:"acronym"`
	AspectRatio string `json:"aspect_ratio"`
	Events      []struct {
		ConferenceURL    string   `json:"conference_url"`
		Date             string   `json:"date"`
		Description      string   `json:"description"`
		FrontendLink     string   `json:"frontend_link"`
		GUID             string   `json:"guid"`
		Length           int64    `json:"length"`
		Link             string   `json:"link"`
		OriginalLanguage string   `json:"original_language"`
		Persons          []string `json:"persons"`
		PosterURL        string   `json:"poster_url"`
		ReleaseDate      string   `json:"release_date"`
		Slug             string   `json:"slug"`
		Subtitle         string   `json:"subtitle"`
		Tags             []string `json:"tags"`
		ThumbURL         string   `json:"thumb_url"`
		Title            string   `json:"title"`
		UpdatedAt        string   `json:"updated_at"`
		URL              string   `json:"url"`
	} `json:"events"`
	ImagesURL      string `json:"images_url"`
	LogoURL        string `json:"logo_url"`
	RecordingsURL  string `json:"recordings_url"`
	ScheduleURL    string `json:"schedule_url"`
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	UpdatedAt      string `json:"updated_at"`
	URL            string `json:"url"`
	WebgenLocation string `json:"webgen_location"`
}

func findEvents(url string) (Event, error) {
	event := Event{}

	r, err := http.Get(url)
	if err != nil {
		return event, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&event)
	return event, err
}
