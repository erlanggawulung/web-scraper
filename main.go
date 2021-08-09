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

type Input struct {
	AllowedDomains []string `json:"allowedDomains"`
	SiteURL        string   `json:"siteURL"`
	ParentClass    string   `json:"parentClass"`
}

func main() {
	input := readInput()
	allSites := make([]Site, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains(input.AllowedDomains...),
	)

	collector.OnHTML(input.ParentClass, func(element *colly.HTMLElement) {
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

	collector.Visit(input.SiteURL)

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

func readInput() Input {
	// Open our jsonFile
	jsonFile, err := os.Open("input.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened input.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var input Input

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &input)

	return input
}
