package scraper

import (
	"strings"

	"github.com/gocolly/colly/v2"
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
func GetSelectorHandlers(pageData *[]CountiesOfTheWorld) []SelectorHandler {
	return []SelectorHandler{
		{
			Name:     "Country",
			Selector: "h3.country-name",
			Handler:  countryName(pageData),
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

func sanatizeString(s string) string {
	// Add any additional sanitization logic as needed
	s = strings.TrimSpace(s)
	return s
}
