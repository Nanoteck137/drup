package cmd

import (
	"log"
	"os"
	"path"
	"slices"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/nanoteck137/swadloon/utils"
	"github.com/spf13/cobra"
)

type Manga struct {
	Name        string
	CoverArtUrl string
	Chapters    []Chapter
}

type Chapter struct {
	Number int
	Name   string
	Url    string
}

type Page struct {
	Url string
}

type Mangapill struct {
	baseUrl string
}

func (m *Mangapill) GetMangaDetails(id string) (Manga, error) {
	res, err := utils.FetchHtml(m.baseUrl + "/manga/" + id)
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

	coverArtUrl := doc.Find(".container").Find("img").AttrOr("data-src", "")

	doc.Find("div[data-filter-list] a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		url, _ := s.Attr("href")

		chapters = append(chapters, Chapter{
			Number: 0,
			Name:   name,
			Url:    url,
		})
	})

	slices.Reverse(chapters)

	for i := range chapters {
		chapters[i].Number = i + 1
	}

	return Manga{
		Name:        name,
		CoverArtUrl: coverArtUrl,
		Chapters:    chapters,
	}, nil
}

func (m *Mangapill) GetChapterPages(chapter Chapter) ([]Page, error) {
	res, err := utils.FetchHtml(m.baseUrl + chapter.Url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
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

func DownloadChapter(p Mangapill, chapter Chapter, out string) ([]string, error) {
	pages, err := p.GetChapterPages(chapter)
	if err != nil {
		return nil, err
	}

	var res []string

	for i, page := range pages {
		pageData, err := utils.FetchPage(p.baseUrl+chapter.Url, page.Url)
		if err != nil {
			return nil, err
		}

		ext, err := utils.GetExtention(pageData.ContentType)
		if err != nil {
			return nil, err
		}

		name := strconv.Itoa(i)
		p := path.Join(out, name+ext)

		err = os.WriteFile(p, pageData.Data, 0644)
		if err != nil {
			return nil, err
		}

		res = append(res, name+ext)
	}

	return res, nil
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
