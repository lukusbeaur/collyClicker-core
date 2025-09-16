package handlers

import (
	"strings"

	"github.com/gocolly/colly/v2"
	scraper "github.com/lukusbeaur/collyclicker-core/scrape"
)

// Example handler for
type CountiesOfTheWorld struct {
	Country string
	/*capital    string
	population string
	area       string*/
	//For practice write your own handlers and selectors for the rest of the elements on the page
	//Check the Input links for the URL
}

// Update the get Seletor handlers for the elements you want to scrape
func GetSelectorHandlers(pageData *[]CountiesOfTheWorld) []scraper.SelectorHandler {
	return []scraper.SelectorHandler{
		{
			Name:     "Country",             //Name of your handler ( EX. Country, Capital, Population, Area ) - used for output file name
			Selector: "h3.country-name",     //the CSS selector to find the element
			Handler:  countryName(pageData), //your function name here: Defines how you extract the data
		},
	}
}

// Define the function for each selector, to find and parse the data from the HTML element
func countryName(data *[]CountiesOfTheWorld) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		country := sanatizeString(e.Text)
		//fmt.Println(country)
		if country == "" {
			return
		}
		*data = append(*data, CountiesOfTheWorld{Country: country})
	}
}

// parse and clean the string when extracting data before storing it
func sanatizeString(s string) string {
	// Add any additional sanitization logic as needed
	s = strings.TrimSpace(s)
	return s
}
