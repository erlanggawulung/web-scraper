package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Site struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	ImageLink   string `json:"imageLink"`
}

func main() {
	allSites := make([]Site, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("webscraper.io"),
	)

	collector.OnHTML(".row", func(element *colly.HTMLElement) {
		log.Println("Element: ", element)
		siteTitle := element.ChildText("h2")

		var siteLink string
		siteLinks := element.ChildAttrs("a", "href")
		if len(siteLinks) > 0 {
			siteLink = siteLinks[0]
		}

		var imageLink string
		imageLinks := element.ChildAttrs("img", "src")
		if len(imageLinks) > 0 {
			imageLink = imageLinks[0]
		}

		siteDescription := element.ChildText("p")

		if len(siteTitle) > 0 {
			site := Site{
				Title:       siteTitle,
				Description: siteDescription,
				Link:        siteLink,
				ImageLink:   imageLink,
			}

			allSites = append(allSites, site)
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://webscraper.io/test-sites")

	writeJSON(allSites)
	writeCSV(allSites)
}

func writeJSON(data []Site) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}
	_ = ioutil.WriteFile("sites.json", file, 0644)
}

func writeCSV(data []Site) {
	rows := [][]string{
		{"Title", "Description", "Link", "ImageLink"},
	}

	for _, site := range data {
		row := []string{site.Title, site.Description, site.Link, site.ImageLink}
		rows = append(rows, row)
	}

	csvfile, err := os.Create("sites.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)

	for _, row := range rows {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()

	csvfile.Close()
}
