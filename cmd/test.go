package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/nanoteck137/swadloon/library"
	"github.com/spf13/cobra"
)

type Manga struct {
	Name     string
	Chapters []Chapter
}

type Chapter struct {
	Index int
	Name  string
	Url   string
}

type Page struct {
	Url string
}

type Mangapill struct {
	baseUrl string
}

func fetchHtml(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("res.StatusCode: %v\n", res.StatusCode)

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		res.Body.Close()
		return nil, errors.New("Expected Content-Type to be text/html")
	}

	return res, nil
}

func (m *Mangapill) GetMangaDetails(id string) (Manga, error) {
	// f, err := os.Open("work/test.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	res, err := fetchHtml(m.baseUrl + "/manga/"+id)
	if err != nil {
		return Manga{}, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var chapters []Chapter

	name := doc.Find("h1").Text()

	doc.Find("div[data-filter-list] a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		url, _ := s.Attr("href")

		chapters = append(chapters, Chapter{
			Index: 0,
			Name:  name,
			Url:   url,
		})
	})

	slices.Reverse(chapters)

	for i := range chapters {
		chapters[i].Index = i
	}

	return Manga{
		Name:     name,
		Chapters: chapters,
	}, nil
}

func (m *Mangapill) GetChapterPages(chapter Chapter) ([]Page, error) {
	f, err := os.Open("work/test2.html")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, err
	}

	var pages []Page

	doc.Find("chapter-page picture img").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("data-src")
		pages = append(pages, Page{
			Url: url,
		})
	})

	return pages, nil
}

func GetExtention(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpeg", nil
	case "image/png":
		return ".png", nil
	}

	return "", fmt.Errorf("Unknown Content-Type: %v", contentType)
}

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"

type PageData struct {
	Data        []byte
	ContentType string
}

func fetchPage(chapterUrl, url string) (PageData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PageData{}, err
	}

	req.Header.Set("Referer", chapterUrl)
	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return PageData{}, err
	}
	defer res.Body.Close()

	fmt.Printf("res.StatusCode: %v\n", res.StatusCode)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PageData{}, err
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		return PageData{}, errors.New("No Content Type")
	}

	return PageData{
		Data:        data,
		ContentType: contentType,
	}, nil
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		lib, err := library.ReadFromDir("/Volumes/media/manga")
		if err != nil {
			log.Fatal(err)
		}

		var id string
		err = huh.NewInput().
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

		if err != nil {
			log.Fatal(err)
		}

		p := Mangapill{
			baseUrl: "https://mangapill.com",
		}

		manga, err := p.GetMangaDetails(id)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(manga)

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
	rootCmd.AddCommand(testCmd)
}
