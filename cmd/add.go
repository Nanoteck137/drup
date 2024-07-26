package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/sewaddle-core/library"
	"github.com/spf13/cobra"
)

func askForId() (string, error) {
	var id string
	err := huh.NewInput().
		Title("Mangapill Id").
		Value(&id).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("can't be empty")
			}

			if _, err := strconv.Atoi(s); err != nil {
				return errors.New("needs to be a number")
			}

			return nil
		}).
		Run()

	return id, err
}

var addCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		lib, err := library.ReadFromDir("./work/lib")
		if err != nil {
			log.Fatal(err)
		}

		p := Mangapill{
			baseUrl: "https://mangapill.com",
		}

		id, err := askForId()
		if err != nil {
			log.Fatal(err)
		}

		manga, err := p.GetMangaDetails(id)
		if err != nil {
			log.Fatal(err)
		}

		var correct bool
		err = huh.NewConfirm().
			Title(fmt.Sprintf("Found '%s'", manga.Name)).
			Affirmative("Correct").
			Negative("Incorrect").
			Value(&correct).
			Run()

		if !correct {
			return
		}

		serie := library.SerieMetadata{
			Title:    manga.Name,
			Chapters: []library.ChapterMetadata{},
			Extra:    map[string]any{},
		}

		serie.Extra["provider"] = "mangapill"
		serie.Extra["mangapill-id"] = id

		err = lib.AddSerie(serie)
		if err != nil {
			// TODO(patrik): Ask the user if they want to change the name
			log.Fatal(err)
		}

		err = lib.FlushToDisk()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
