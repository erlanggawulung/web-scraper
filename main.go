package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gocolly/colly"
)

type Site struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

func main() {
	allSites := make([]Site, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("webscraper.io"),
	)

	collector.OnHTML(".col-md-7", func(element *colly.HTMLElement) {
		log.Println("Element: ", element)
		siteTitle := element.ChildText("h2")
		siteLink := element.ChildAttrs("a", "href")[0]
		log.Println("siteLink: ", siteLink)
		siteDescription := element.ChildText("p")
		site := Site{
			Title:       siteTitle,
			Description: siteDescription,
			Link:        siteLink,
		}

		allSites = append(allSites, site)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://webscraper.io/test-sites")

	writeJSON(allSites)
}

func writeJSON(data []Site) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}
	_ = ioutil.WriteFile("sites.json", file, 0644)
}
