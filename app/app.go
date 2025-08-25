package app

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/lukusbeaur/collyclicker-core/internal/Util"
	"github.com/lukusbeaur/collyclicker-core/internal/csvparser"
	"github.com/lukusbeaur/collyclicker-core/internal/fileutils"
)

func Run() error {
	fmt.Println("CollyClicker Application Running")

	//--------------Create Colly Collector ----------------
	// -------------Non Call back functions ----------------
	c := colly.NewCollector(
		colly.AllowedDomains("scrapethissite.com"),
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
		r.Headers.Set("Referer", "https://fbref.com/")
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

	}
	return nil
}
