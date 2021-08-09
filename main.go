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

type Input struct {
	AllowedDomains []string   `json:"allowedDomains"`
	SiteURL        string     `json:"siteURL"`
	ParentClass    string     `json:"parentClass"`
	InputMaps      []InputMap `json:"maps"`
	JSONFileName   string     `json:"jsonFileName"`
	CSVFileName    string     `json:"csvFileName"`
}

type InputMap struct {
	Key        string   `json:"key"`
	ChildText  string   `json:"childText"`
	ChildAttrs []string `json:"childAttrs"`
}

type OutputMap struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	input := readInput()
	results := make([][]OutputMap, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains(input.AllowedDomains...),
	)

	collector.OnHTML(input.ParentClass, func(element *colly.HTMLElement) {
		log.Println("Element: ", element)
		outputRow := make([]OutputMap, 0)
		for _, item := range input.InputMaps {
			if len(item.ChildText) > 0 {
				outputItem := OutputMap{
					Key:   item.Key,
					Value: element.ChildText(item.ChildText),
				}
				outputRow = append(outputRow, outputItem)
			} else if len(item.ChildAttrs) > 0 {
				values := element.ChildAttrs(item.ChildAttrs[0], item.ChildAttrs[1])
				if len(values) > 0 {
					outputItem := OutputMap{
						Key:   item.Key,
						Value: values[0],
					}
					outputRow = append(outputRow, outputItem)
				}
			}
		}
		results = append(results, outputRow)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit(input.SiteURL)

	if len(input.JSONFileName) > 0 {
		writeJSON(input, results)
	}
	if len(input.CSVFileName) > 0 {
		writeCSV(input, results)
	}
}

func writeJSON(input Input, data [][]OutputMap) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}
	_ = ioutil.WriteFile("output/"+input.JSONFileName, file, 0644)
}

func writeCSV(input Input, data [][]OutputMap) {
	rows := [][]string{}

	// Generate header, rows[0]
	rowZero := []string{}
	for _, item := range input.InputMaps {
		rowZero = append(rowZero, item.Key)
	}
	rows = append(rows, rowZero)

	for _, row := range data {
		rowN := []string{}
		if len(row) == len(rowZero) {
			for _, column := range row {
				rowN = append(rowN, column.Value)
			}
			rows = append(rows, rowN)
		}
	}

	csvfile, err := os.Create("output/" + input.CSVFileName)

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
