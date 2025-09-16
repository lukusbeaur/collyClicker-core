package app

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gocolly/colly/v2"
	appparams "github.com/lukusbeaur/collyclicker-core/app/config"
	handlers "github.com/lukusbeaur/collyclicker-core/app/handlers"
	"github.com/lukusbeaur/collyclicker-core/internal/Util"
	"github.com/lukusbeaur/collyclicker-core/internal/csvparser"
	"github.com/lukusbeaur/collyclicker-core/internal/fileutils"
	scraper "github.com/lukusbeaur/collyclicker-core/scrape"
)

func Run() error {
	fmt.Println("CollyClicker Application Running")
	//Call the default config function to set the default values
	cfg := appparams.DefaultConfig()

	// -------------handler selection -----------------
	var pageData []handlers.CountiesOfTheWorld
	selectorList := handlers.GetSelectorHandlers(&pageData)

	// -------------Create Colly Collector -----------------
	c := scraper.NewCollectorFromAppConfig(cfg, selectorList, nil)

	//--------------Call back functions ----------------
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://scrapethissite.com")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})
	//temporary debugging on failed requests. once I add logging back this will change
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s\n", r.Request.URL)
		fmt.Printf("Status Code: %d\n", r.StatusCode)
		fmt.Printf("Error: %v\n", err)
	})

	c.OnScraped(func(r *colly.Response) {
		for _, it := range pageData {
			_ = fileutils.WriteLineCSV("countries.csv", cfg.OutputDir, []string{it.Country})
		}
	})
	//--------------Colly setup Complete ----------------

	//--------------CSV of URLS to scrape ----------------
	files, err := fileutils.Findcsvfiles(cfg.InputDir)
	if err != nil {
		Util.Logger.Error("Trouble finding csvs in directory",
			"Location", "app.go - FindcsvFiles",
			"Error", err)
	}
	// _ = index of csvArray if needed change to for index, record := range csvArray
	for _, record := range files {
		fmt.Println("Processing file:", record)
		//If you want to ping the URLs first, uncomment the next line
		//Util.CheckURL(record)

		iter, err := csvparser.NewCSViter(cfg.InputDir, record)
		if err != nil {
			Util.Logger.Error("Trouble opening CSV file and or Iterator",
				"Location", "app.go - Range csvArray loop",
				"Record", record,
				"Error", err)
		}
		defer iter.Close()

		for {
			//pageData - creates a slice of county structs to hold the data
			pageData := []handlers.CountiesOfTheWorld{}
			row, _, _, err := iter.Next()
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
			pageData = pageData[:0]
			if err := c.Visit(row[0]); err != nil {
				fmt.Println("Visit:", err)
				continue
			}
		}
	}
	c.Wait()
	return nil
}
