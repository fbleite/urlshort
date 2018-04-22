package urlshort

import (
	"testing"
	"log"
)

func Test_parseYAML(t *testing.T) {
	yamlInput := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlPathToUrlArray, err:= parseYAML([]byte(yamlInput))

	if err != nil {
		t.Fatalf("Error should be nil")
	}

	if (yamlPathToUrlArray[0].Path != "/urlshort" ||
		yamlPathToUrlArray[0].URL != "https://github.com/gophercises/urlshort"||
		yamlPathToUrlArray[1].Path != "/urlshort-final" ||
		yamlPathToUrlArray[1].URL != "https://github.com/gophercises/urlshort/tree/solution"){

			t.Fail()
	}
}


func Test_parseJSON(t *testing.T) {
	jsonInput := `[{"path": "/this-path", "url": "http://follow-the-path.org"},
{"path": "/that-path", "url": "http://follow-that-path.org"}]`
	parsedJson, err := parseJSON([]byte(jsonInput))
	if err != nil{
		t.Fatalf(err.Error())
	}

	log.Println(parsedJson)

}

func Test_accessBoltDb(t *testing.T) {
	accessBoltDb()
}
