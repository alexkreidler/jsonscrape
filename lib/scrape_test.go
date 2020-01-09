package lib

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
	"runtime"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary



var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)


func TestScrape(t *testing.T) {
	data, err := ioutil.ReadFile("./test/scrape_rcp.json")
	if err != nil {
		t.Fatal(err)
	}

	parsed := Config{}
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Error("failed to parse json", err)
	}

	s, err := NewScraper(parsed)
	if err != nil {
		t.Error("failed to initialize scraper", err)
	}

	d, e := s.Scrape()
	if e != nil {
		t.Error("failed to get results", e)
	}

	b, err := json.MarshalIndent(d, "", "    ")
	t.Log(string(b))

	filepath := path.Join(basepath, "./test/scrape_output.json")
	t.Log(filepath)
	err = ioutil.WriteFile(filepath, b, 0644)
	if err != nil {
		t.Error("failed to write output file", err)
	}
}
