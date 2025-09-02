package app

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	cfg "github.com/lukusbeaur/collyclicker-core/app/config"
	handlers "github.com/lukusbeaur/collyclicker-core/app/handlers"
	"github.com/lukusbeaur/collyclicker-core/internal/Util"
	"github.com/lukusbeaur/collyclicker-core/internal/csvparser"
	"github.com/lukusbeaur/collyclicker-core/internal/fileutils"
)

func Run() error {
	fmt.Println("CollyClicker Application Running")
	//Call the default config function to set the default values
	cfg := cfg.ConfigDefaults()
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
	csvArray, err := fileutils.Findcsvfiles("") //defaults to Input_Links/
	if err != nil {
		Util.Logger.Error("Trouble finding csvs in directory",
			"Location", "app.go - FindcsvFiles",
			"Error", err)
	}
	// _ = index of csvArray if needed change to for index, record := range csvArray
	for _, record := range csvArray {
		fmt.Println("Processing file:", record)
		//If you want to ping the URLs first, uncomment the next line
		//Util.CheckURL(record)

		csvfile, err := csvparser.NewCSViter("", record) //defaults to Input_Links/links.csv if both "".
		if err != nil {
			Util.Logger.Error("Trouble opening CSV file and or Iterator",
				"Location", "app.go - Range csvArray loop",
				"Record", record,
				"Error", err)
		}
		defer csvfile.Close()

		for {
			//pageData - creates a slice of county structs to hold the data
			pageData := []handlers.CountiesOfTheWorld{}
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
			for _, h := range handlers.GetSelectorHandlers(&pageData) {
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
