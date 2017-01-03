package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"
)

const (
	noTranslations = true
	downloadPath   = "./downloads/"
)

var (
	preferredMimeTypes    = []string{"video/webm", "video/mp4", "video/ogg", "audio/ogg", "audio/opus", "audio/mpeg", "application/x-subrip"}
	extensionForMimeTypes = make(map[string]string)
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

func priorityForMimeType(mime string) int {
	for i, v := range preferredMimeTypes {
		if strings.ToLower(mime) == v {
			return i
		}
	}

	return -1
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
	extensionForMimeTypes["video/webm"] = "webm"
	extensionForMimeTypes["video/mp4"] = "mp4"
	extensionForMimeTypes["video/ogg"] = "ogm"
	extensionForMimeTypes["audio/ogg"] = "ogg"
	extensionForMimeTypes["audio/opus"] = "opus"
	extensionForMimeTypes["audio/mpeg"] = "mp3"

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

			bestMatch := Recording{}
			highestPriority := -1
			if len(media.Recordings) == 0 {
				panic("No recordings found for this event!")
			}
			for _, m := range media.Recordings {
				if m.Width == 0 {
					fmt.Printf("\t\tFound audio (%s): %d minutes (HD: %t, %dMiB) %s\n", m.MimeType, m.Length/60, m.HighQuality, m.Size, m.URL)
				} else {
					fmt.Printf("\t\tFound video (%s): %d minutes, %dx%d (HD: %t, %dMiB) %s\n", m.MimeType, m.Length/60, m.Width, m.Height, m.HighQuality, m.Size, m.URL)
				}

				prio := priorityForMimeType(m.MimeType)
				if prio < 0 {
					panic("Unknown mimetype encountered:" + m.MimeType)
				}

				if prio < highestPriority || (prio == highestPriority && m.Width > bestMatch.Width) || highestPriority == -1 {
					highestPriority = prio
					bestMatch = m
				}
			}

			path := filepath.Join(downloadPath, sanitize.Path(v.Title))
			basename := sanitize.BaseName(e.Title) + "." + extensionForMimeTypes[bestMatch.MimeType]
			filename := filepath.Join(path, basename)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				os.MkdirAll(path, 0755)

				fmt.Println("\t\tDownloading:", bestMatch.RecordingURL)
				out, err := os.Create(filename)
				if err != nil {
					panic(err)
				}
				defer out.Close()

				resp, err := http.Get(bestMatch.RecordingURL)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				_, err = io.Copy(out, resp.Body)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("\t\tFile", filename, "already exists - skipping!")
			}

			fmt.Println()
		}
	}

	fmt.Println("Done.")
}
