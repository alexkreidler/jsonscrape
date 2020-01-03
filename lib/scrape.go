package lib

import (
	// "encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/imdario/mergo"
)

// Scraper represents a web scraper
type Scraper struct {
	config    JSONSpec
	collector *colly.Collector
}

// JSONSpec represents the JSON file specification of the scraper. It includes general settings, colly specific settings, and the raw scraper spec
type JSONSpec struct {
	GeneralConfig `json:"config"`
	CollyConfig   colly.Collector `json:"colly""`
	ScrapeConfig  `json:"scraper"`
}

// GeneralConfig contains the general configuration for jsonscrape
type GeneralConfig struct {
	Sites []string
}

// JSON represents a JSON input
type JSON []byte

// NewScraper returns a new JSONscraper
func NewScraper(parsed JSONSpec) (Scraper, error) {
	c := &colly.Collector{}
	c.Init()

	// merge the JSON colly spec into the real colly collector
	if err := mergo.Merge(c, &parsed.CollyConfig, mergo.WithOverride); err != nil {
		return Scraper{}, err
	}

	return Scraper{
		config:    parsed,
		collector: c,
	}, nil
}

// Value represents the available values to access from the selected elements
type Value string

func updateData(d *Results, scrapeValue Value, scrapeName string) func(e *colly.HTMLElement) {
	(*d)[scrapeName] = make(chan Value)
	return func(e *colly.HTMLElement) {
		var val Value

		z := func() {
			switch scrapeValue {
			case "text":
				val = Value(strings.Trim(e.Text, "\n "))
				break
			default:
				val = Value(strings.Trim(e.Attr(string(scrapeValue)), "\n "))
			}
			// fmt.Println("got: ", val, scrapeName)
			// naive way: just append to the slice. this may not be concurrency safe
			// in the future, Result will proably be a channel allowing different pages on the site which are accessed
			// at a later time to add more and more matches for the same selector
			cur := (*d)[scrapeName]
			if cur != nil {
				cur <- val
			}
		}

		go z()

	}

}

// ScrapeConfig represents the values that the scraper will retrieve from the various sites
type ScrapeConfig map[string]datum

type datum struct {
	Selector string `json:"selector,omitempty"`
	Value    `json:"value,omitempty"`
}

type Results map[string]Result

// type Result []Value

type Result chan Value

// Scrape runs the scraper as specified. It returns the data retrieved and/or an error from the scraping process
func (s *Scraper) Scrape() (interface{}, error) {
	c := *s.collector

	fmt.Printf("%+#v\n", s.config.ScrapeConfig)

	results := Results{} //make(Results)

	// k and v are pass by reference in loops, except across function calls, so we extract to the updateData function

	// fmt.Printf("")
	for k, v := range s.config.ScrapeConfig {
		fmt.Println("setting up:", k, v)
		c.OnHTML(v.Selector, updateData(&results, v.Value, k))
	}

	// Find and visit all links
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	fmt.Println("found link")
	// 	// _ = e.Request.Visit(e.Attr("href"))
	// })
	// c.OnHTML("h2", func(e *colly.HTMLElement) {
	// 	fmt.Println("found header")
	// 	fmt.Println(strings.Trim(e.Text, "\n "))
	// })
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// 	fmt.Println(r.Body)
	// })
	c.OnResponse(func(r *colly.Response) {
		// fmt.Println(string(r.Body))
	})

	for _, site := range s.config.GeneralConfig.Sites {
		fmt.Println("\nsetting up site:", site)
		if err := c.Visit(site); err != nil {
			return nil, err
		}
	}

	
	fmt.Println("done")
	// close(results["race_url"])
	for _, item := range <-results["race_url"] {
		fmt.Println(item)
	}

	return nil, nil
}
