package library

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/flytam/filenamify"
)

type ChapterMetadata struct {
	Index int      `json:"index"`
	Name  string   `json:"name"`
	Pages []string `json:"pages"`
}

type SerieMetadata struct {
	Title    string            `json:"title"`
	Chapters []ChapterMetadata `json:"chapters"`
	Extra    map[string]any    `json:"extra,omitempty"`

	new bool
}

type Library struct {
	Base   string
	Series []SerieMetadata
}

func ReadFromDir(dir string) (*Library, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var series []SerieMetadata

	for _, entry := range entries {
		// Skip over files and folders starting with a dot
		if entry.Name()[0] == '.' {
			continue
		}

		p := path.Join(dir, entry.Name())

		data, err := os.ReadFile(path.Join(p, "manga.json"))
		if err != nil {
			log.Printf("Warning: '%v' has no manga.json", p)
			return nil, err
		}

		var metadata SerieMetadata
		err = json.Unmarshal(data, &metadata)
		if err != nil {
			return nil, err
		}

		series = append(series, metadata)
	}

	// pretty.Println(series)

	return &Library{
		Base:   dir,
		Series: series,
	}, nil
}

func (lib *Library) AddSerie(serie SerieMetadata) error {
	for _, sm := range lib.Series {
		if sm.Title == serie.Title {
			return fmt.Errorf("Serie with name '%v' already exists", serie.Title)
		}
	}

	serie.new = true
	lib.Series = append(lib.Series, serie)
	return nil
}

func (lib *Library) FlushToDisk() error {
	for _, serie := range lib.Series {
		if serie.new {
			title, err := filenamify.FilenamifyV2(serie.Title, func(options *filenamify.Options) {
				options.Replacement = ""
			})
			if err != nil {
				return err
			}

			d := path.Join(lib.Base, title)
			err = os.Mkdir(d, 0755)
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(serie, "", "  ")
			if err != nil {
				return err
			}

			out := path.Join(d, "manga.json")
			err = os.WriteFile(out, data, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
