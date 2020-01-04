package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/kennygrant/sanitize"
)

const ()

var (
	preferredMimeTypes    = []string{"video/webm", "video/mp4", "video/ogg", "audio/ogg", "audio/opus", "audio/mpeg", "application/x-subrip"}
	extensionForMimeTypes = make(map[string]string)

	downloadPath string
	name         string
	language     string
	listOnly     bool
)

type Conference struct {
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
}

type Conferences struct {
	Conferences []Conference `json:"conferences"`
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
	flag.StringVar(&name, "name", "", "download media of a specific conference only (e.g. '33c3')")
	flag.StringVar(&downloadPath, "destination", "./downloads/", "where to store downloaded media")
	flag.StringVar(&language, "language", "", "preferred language if available (eng, deu or fra)")
	flag.Parse()

	if len(flag.Args()) > 0 {
		arg := strings.ToLower(flag.Args()[0])
		listOnly = arg == "list"
	}

	name = strings.ToLower(name)
	language = strings.ToLower(language)

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

	found := false
	for _, v := range ci.Conferences {
		if !listOnly && len(name) > 0 && name != strings.ToLower(v.Acronym) {
			continue
		}
		fmt.Printf("Conference: %s (%s)\n", v.Acronym, v.Title)
		found = true
		if listOnly {
			continue
		}

		events, err := findEvents(v.URL)
		if err != nil {
			panic(err)
		}

		for _, e := range events.Events {
			desc := strings.Replace(sanitize.HTML(e.Description), "\n", "", -1)
			if len(desc) > 48 {
				desc = desc[:45] + "..."
			}
			if len(desc) > 0 {
				desc = " - " + desc
			}
			fmt.Printf("Event: %s%s\n", e.Title, desc)

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
				if (len(language) == 0 || language != strings.ToLower(m.Language)) &&
					m.Language != e.OriginalLanguage {
					continue
				}

				if m.Width == 0 {
					fmt.Printf("\tFound other/audio (%s): %d minutes (HD: %t, %dMiB) %s\n", m.MimeType, m.Length/60, m.HighQuality, m.Size, m.URL)
				} else {
					fmt.Printf("\tFound video (%s): %d minutes, %dx%d (HD: %t, %dMiB) %s\n", m.MimeType, m.Length/60, m.Width, m.Height, m.HighQuality, m.Size, m.URL)
				}

				prio := priorityForMimeType(m.MimeType)
				if prio < 0 {
					fmt.Println("Unknown mimetype encountered:" + m.MimeType)
				}

				pick := false
				if highestPriority == -1 {
					// if we got nothing so far, always pick any available option
					pick = true
				}
				if strings.ToLower(bestMatch.Language) != language && strings.ToLower(m.Language) == language {
					// we already found something, but this is the preferred language
					pick = true
				} else {
					if prio < highestPriority || (prio == highestPriority && m.Width > bestMatch.Width) {
						// we already found something, but this has a higher resolution
						pick = true
					}
				}

				if pick {
					highestPriority = prio
					bestMatch = m
				}
			}

			if len(bestMatch.RecordingURL) == 0 {
				fmt.Println("Could not find any desired version of this event, sorry. Skipping!")
				continue
			}
			err = download(v, e, bestMatch)
			if err != nil {
				panic(err)
			}

			fmt.Println()
		}
	}

	if found {
		fmt.Println("Done.")
	} else {
		fmt.Println("Couldn't find any conference with acronym", name)
	}
}
