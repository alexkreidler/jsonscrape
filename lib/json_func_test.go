package lib

import (
	"github.com/davecgh/go-spew/spew"
	colly "github.com/gocolly/colly/v2"
	"testing"
)

// This test file verifies that we can successfully marshal the Colly struct to JSON
// For structs that contain functions, the function should either be private, like in Colly v2, or it should have a `json:"-"` tag so it is ignored by the marshaller

type TestComplex struct {
	ObjName   string
	Age       int
	MyHandler func(path string) string `json:"-"`
}

func TestJSONFun(t *testing.T) {
	o := TestComplex{
		ObjName: "johhny be good",
		Age:     245,
		MyHandler: func(path string) string {
			return "hello world: " + path
		},
	}
	spew.Dump(o)

	res, err := json.MarshalIndent(&o, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestCollyMarshal(t *testing.T) {
	c := colly.NewCollector()

	spewConfig := spew.ConfigState{DisableMethods: true, DisablePointerMethods: true}
	spewConfig.Dump(c)

	res, err := json.MarshalIndent(&c, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(res))
}
