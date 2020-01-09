package lib

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/imdario/mergo"
)

// Scraper represents a web scraper
type Scraper struct {
	Config    Config
	Collector *colly.Collector
	Results   Results
}

// Config represents the JSON file specification of the scraper. It includes general settings, colly specific settings, and the raw scraper spec
type Config struct {
	GeneralConfig `json:"config"`
	CollyConfig   colly.Collector `json:"colly"`
	ScrapeConfig  `json:"scraper"`
}

// GeneralConfig contains the general configuration for jsonscrape
type GeneralConfig struct {
	Sites []string
}

// JSON represents a JSON input
type JSON []byte

// NewScraper returns a new JSONscraper
func NewScraper(config Config) (Scraper, error) {
	c := &colly.Collector{}
	c.Init()

	// merge the JSON colly spec into the real colly Collector
	if err := mergo.Merge(c, &config.CollyConfig, mergo.WithOverride); err != nil {
		return Scraper{}, err
	}

	return Scraper{
		Config:    config,
		Collector: c,
	}, nil
}

// Maps are passed by reference
func updateData(d Results, scrapeValues ValueMap, scrapeName string) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		vm := ValueMap{}

		for outName, attrName := range scrapeValues {
			var val Value
			switch attrName {
			case "text":
				val = Value(strings.Trim(e.Text, "\n "))
				break
			default:
				val = Value(strings.Trim(e.Attr(string(attrName)), "\n "))
			}

			vm[outName] = val
		}
		// naive way: just append to the slice. this may not be concurrency safe
		// in the future, Result will proably be a channel allowing different pages on the site which are accessed
		// at a later time to add more and more matches for the same selector
		cur, ok := d[scrapeName]
		if !ok {
			d[scrapeName] = []ValueMap{vm}
		} else {
			d[scrapeName] = append(cur, vm)
		}
	}

}

// ScrapeConfig represents the values that the scraper will retrieve from the various sites
type ScrapeConfig map[string]datum

// A datum represents a single element selected via a selector
type datum struct {
	Selector string `json:"selector,omitempty"`
	Values ValueMap `json:"values,omitempty"`
}

// Value represents the available values to access from the selected elements
// It corresponds to the attributes on a given element
type Value string

// A value map is a map from strings or keys to values which will be extracted from the given element
type ValueMap map[string]Value

// Results is a map of values that are the results of the scraper
// Uses a pointer to allow nil base types
type Results map[string][]ValueMap

// type Result []Value

//type Result

// Scrape runs the scraper as specified. It returns the data retrieved and/or an error from the scraping process
// It also will block until the scraper has stopped entirely
func (s *Scraper) Scrape() (interface{}, error) {
	c := s.Collector

	results := Results{}

	// k and v are pass by reference in loops, except across function calls, so we extract to the updateData function
	for k, v := range s.Config.ScrapeConfig {
		//fmt.Println("setting up:", k, v)
		c.OnHTML(v.Selector, updateData(results, v.Values, k))
	}

	for _, site := range s.Config.GeneralConfig.Sites {
		//fmt.Println("\nsetting up site:", site)
		if err := c.Visit(site); err != nil {
			return nil, err
		}
	}

	c.Wait()

	return results, nil
}
