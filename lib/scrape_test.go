package lib

import (
	"io/ioutil"
	"testing"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func TestScrape(t *testing.T) {
	data, err := ioutil.ReadFile("./test/scrape_rcp.json")
	if err != nil {
		t.Fatal(err)
	}

	parsed := JSONSpec{}
	err = json5.Unmarshal(data, &parsed)
	if err != nil {
		t.Error("failed to parse json", err)
	}

	s, err := NewScraper(parsed)
	if err != nil {
		t.Error("failed to initialize scraper", err)
	}

	d, e := s.Scrape()
	t.Log(d)
	if e != nil {
		t.Error("failed to initialize scraper", e)
	}
}
