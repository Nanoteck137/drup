package cmd

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/kr/pretty"
	"github.com/nanoteck137/sewaddle/library"
	"github.com/nanoteck137/drup/utils"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		lib, err := library.ReadFromDir("./work/lib")
		if err != nil {
			log.Fatal(err)
		}

		p := Mangapill{
			baseUrl: "https://mangapill.com",
		}

		for i := range lib.Series {
			serie := &lib.Series[i]
			provider, ok := serie.Extra["provider"]
			if !ok {
				continue
			}

			ps := provider.(string)

			switch ps {
			case "mangapill":
				id, ok := serie.Extra["mangapill-id"]
				if !ok {
					log.Fatal("No id set")
				}

				m, err := p.GetMangaDetails(id.(string))
				if err != nil {
					log.Fatal(err)
				}

				pretty.Println(m)

				if m.CoverArtUrl != "" && serie.CoverArt == "" {
					data, err := utils.FetchPage(p.baseUrl, m.CoverArtUrl)
					if err != nil {
						log.Fatal(err)
					}

					ext, err := utils.GetExtention(data.ContentType)
					if err != nil {
						log.Fatal(err)
					}

					name := "cover" + ext
					p := path.Join(serie.Path(), name)
					err = os.WriteFile(p, data.Data, 0644)
					if err != nil {
						log.Fatal(err)
					}

					serie.CoverArt = name
					serie.MarkChanged()
				}

				var missing []int
				for i, c := range m.Chapters {
					found := false
					for _, serieChapter := range serie.Chapters {
						if c.Number == serieChapter.Number {
							found = true
							break
						}
					}

					if found {
						continue
					}

					missing = append(missing, i)
				}

				pretty.Println(missing)

				chaptersPath := path.Join(serie.Path(), "chapters")
				err = os.Mkdir(chaptersPath, 0755)
				if err != nil {
					if !os.IsExist(err) {
						log.Fatal(err)
					}
				}

				for _, i := range missing {
					chapter := m.Chapters[i]
					pretty.Println(chapter)
					out := path.Join(serie.Path(), "chapters", strconv.Itoa(chapter.Number))

					err := os.Mkdir(out, 0755)
					if err != nil {
						if os.IsExist(err) {
							err := os.RemoveAll(out)
							if err != nil {
								log.Fatal(err)
							}

							err = os.Mkdir(out, 0755)
							if err != nil {
								log.Fatal(err)
							}
						} else {
							log.Fatal(err)
						}
					}

					pages, err := DownloadChapter(p, chapter, out)
					if err != nil {
						log.Fatal(err)
					}

					pretty.Println(pages)
					serie.Chapters = append(serie.Chapters, library.ChapterMetadata{
						Number: chapter.Number,
						Name:   chapter.Name,
						Pages:  pages,
					})

					serie.MarkChanged()

					break
				}
			}

			lib.FlushToDisk()
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
