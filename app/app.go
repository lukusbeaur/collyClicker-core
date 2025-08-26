package app

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/lukusbeaur/collyclicker-core/internal/Util"
	"github.com/lukusbeaur/collyclicker-core/internal/csvparser"
	"github.com/lukusbeaur/collyclicker-core/internal/fileutils"
	scraper "github.com/lukusbeaur/collyclicker-core/scrape"
)

func Run() error {
	fmt.Println("CollyClicker Application Running")

	//--------------Create Colly Collector ----------------
	// -------------Non Call back functions ----------------
	c := colly.NewCollector(
		colly.AllowedDomains("scrapethissite.com", "www.scrapethissite.com"),
		colly.Async(false),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       2 * time.Second, // min
		RandomDelay: 4 * time.Second, // extra random
	})
	// Dont ignore robot.txt but allow domain revisiting
	c.IgnoreRobotsTxt = false
	c.AllowURLRevisit = true
	//--------------Call back functions ----------------
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://scrapethissite.com")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})
	//temporary debugging on failed requests.
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s\n", r.Request.URL)
		fmt.Printf("Status Code: %d\n", r.StatusCode)
		fmt.Printf("Error: %v\n", err)
	})

	//--------------Colly setup Complete ----------------

	//--------------CSV of URLS to scrape ----------------
	dirname_ready := "Input_Links/"
	csvArray, err := fileutils.Findcsvfiles(dirname_ready)
	if err != nil {
		Util.Logger.Error("Trouble finding csvs in dirname_ready",
			"Location", "app.go - FindcsvFiles",
			"dirname_ready", dirname_ready,
			"Error", err)
	}
	// _ = index of csvArray if needed change to for index, record := range csvArray
	for _, record := range csvArray {
		fmt.Println("Processing file:", record)
		//If you want to ping the URLs first, uncomment the next line
		//Util.CheckURL(record)

		csvfile, err := csvparser.NewCSViter(dirname_ready + record)
		if err != nil {
			Util.Logger.Error("Trouble opening CSV file and or Iterator",
				"Location", "app.go - Range csvArray loop",
				"Record", record,
				"Error", err)
		}
		defer csvfile.Close()

		for {
			//pageData - creates a slice of county structs to hold the data
			pageData := []scraper.CountiesOfTheWorld{}
			row, _, _, err := csvfile.Next()
			if errors.Is(err, io.EOF) {
				fmt.Println("End of CSV file reached")
				break
			}
			if err != nil {
				fmt.Println("Error reacding row:", err)
				continue
			}
			if len(row) == 0 || !strings.HasPrefix(row[0], "http") {
				continue
			}
			// this gets the first element of the row which is the URL
			url := row[0]
			for _, h := range scraper.GetSelectorHandlers(&pageData) {
				//on HTML calls the handler function for each selector
				c.OnHTML(h.Selector, func(e *colly.HTMLElement) { h.Handler(e) })
				c.OnScraped(func(r *colly.Response) {
					fmt.Println("Finished", r.Request.URL)

					for _, it := range pageData {
						fmt.Printf("Country: %s\n", it.Country)
						//The example uses h.name which is the handler name you will assign in /scrape/handlers.go
						//this will default and save to Output_Data/
						fileutils.WriteLineCSV(h.Name+".csv", []string{it.Country})
					}

				})
			}
			if err := c.Visit(url); err != nil {
				fmt.Println("Error visiting URL:", err)
				continue
			}
		}
		c.Wait()
	}
	return nil
}
