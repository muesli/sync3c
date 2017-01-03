package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kennygrant/sanitize"
)

type Conferences struct {
	Conferences []struct {
		Acronym        string `json:"acronym"`
		AspectRatio    string `json:"aspect_ratio"`
		ImagesURL      string `json:"images_url"`
		LogoURL        string `json:"logo_url"`
		RecordingsURL  string `json:"recordings_url"`
		ScheduleURL    string `json:"schedule_url"`
		Slug           string `json:"slug"`
		Title          string `json:"title"`
		UpdatedAt      string `json:"updated_at"`
		URL            string `json:"url"`
		WebgenLocation string `json:"webgen_location"`
	} `json:"conferences"`
}

func findConferences(url string) (Conferences, error) {
	ci := Conferences{}

	r, err := http.Get(url)
	if err != nil {
		return ci, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&ci)
	return ci, err
}

func main() {
	ci, err := findConferences("https://api.media.ccc.de/public/conferences")
	if err != nil {
		panic(err)
	}

	for _, v := range ci.Conferences {
		fmt.Printf("Found conference: %s, URL: %s\n", v.Title, v.URL)

		events, err := findEvents(v.URL)
		if err != nil {
			panic(err)
		}

		for _, e := range events.Events {
			desc := strings.Replace(sanitize.HTML(e.Description), "\n", "", -1)
			if len(desc) > 48 {
				desc = desc[:45] + "..."
			}
			fmt.Printf("\tFound event: %s - %s\n", e.Title, desc)

			media, err := findMedia(e.URL)
			if err != nil {
				panic(err)
			}

			for _, m := range media.Recordings {
				if m.Width == 0 {
					fmt.Printf("\t\tFound audio (%s): %d minutes (HD: %t, %dMiB)- %s\n", m.MimeType, m.Length/60, m.HighQuality, m.Size, m.URL)
				} else {
					fmt.Printf("\t\tFound video (%s): %d minutes, %dx%d (HD: %t, %dMiB)- %s\n", m.MimeType, m.Length/60, m.Width, m.Height, m.HighQuality, m.Size, m.URL)
				}
			}
		}
	}
}
