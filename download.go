package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"
	"github.com/muesli/goprogressbar"
)

type WriteProgressBar struct {
	ProgressBar *goprogressbar.ProgressBar
}

func SizeToString(size uint64) (str string) {
	b := float64(size)

	switch {
	case size >= 1<<60:
		str = fmt.Sprintf("%.2f EiB", b/(1<<60))
	case size >= 1<<50:
		str = fmt.Sprintf("%.2f PiB", b/(1<<50))
	case size >= 1<<40:
		str = fmt.Sprintf("%.2f TiB", b/(1<<40))
	case size >= 1<<30:
		str = fmt.Sprintf("%.2f GiB", b/(1<<30))
	case size >= 1<<20:
		str = fmt.Sprintf("%.2f MiB", b/(1<<20))
	case size >= 1<<10:
		str = fmt.Sprintf("%.2f KiB", b/(1<<10))
	default:
		str = fmt.Sprintf("%dB", size)
	}

	return
}

func (wc *WriteProgressBar) Write(p []byte) (int, error) {
	n := len(p)
	wc.ProgressBar.Current += int64(n)
	wc.ProgressBar.LazyPrint()

	return n, nil
}

func download(v Conference, e Event, m Recording) error {
	author := ""
	subtitle := ""
	lang := ""
	if len(e.Persons) > 0 {
		author = sanitize.BaseName(e.Persons[0]) + " - "
	}
	if len(e.Subtitle) > 0 {
		subtitle = " (" + sanitize.BaseName(e.Subtitle) + ")"
	}
	if e.OriginalLanguage != m.Language {
		lang = " [" + m.Language + "]"
	}

	path := filepath.Join(downloadPath, sanitize.Path(v.Title))
	basename := fmt.Sprintf("%s%s%s%s", author, sanitize.BaseName(e.Title), subtitle, lang) + "." + extensionForMimeTypes[m.MimeType]
	filename := filepath.Join(path, basename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)

		fmt.Println("Downloading:", m.RecordingURL)
		out, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer out.Close()

		resp, err := http.Get(m.RecordingURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		pb := &goprogressbar.ProgressBar{
			Text:  filename,
			Total: resp.ContentLength,
			Width: 60,
			PrependTextFunc: func(p *goprogressbar.ProgressBar) string {
				return fmt.Sprintf("%s / %s",
					SizeToString(uint64(p.Current)),
					SizeToString(uint64(p.Total)))
			},
		}

		src := io.TeeReader(resp.Body, &WriteProgressBar{ProgressBar: pb})
		_, err = io.Copy(out, src)
		if err != nil {
			return err
		}

		fmt.Println()
	} else {
		fmt.Println("File", filename, "already exists - skipping!")
	}

	return nil
}
