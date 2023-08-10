package cellcube_test

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zewolfe/hermes/pkg/services/cellcube"
)

var (
	menu cellcube.Menu = cellcube.Menu{Desc: "News",
		Pages: []cellcube.Page{
			{
				Body: "Headlines",
				PageLinks: []cellcube.PageLink{
					{
						HREF: "#item1",
						A:    "Interest rates cut",
					},
					{
						HREF: "#item2",
						A:    "Concorde resumes service",
					},
				},
			},
			{
				Tag: "item1",
				Body: `
    WASHINGTON-In a much anticipated move, the Federal Reserve
    announced new rate cuts amid growing economic concerns.`,
				PageLinks: []cellcube.PageLink{
					{
						HREF: "#item2",
						A:    "Next article",
					},
				},
			},
			{
				Tag: "item2",
				Body: `
    PARIS-Air France resumed its Concorde service Monday.
    The plane had been grounded following a tragic accident.
    `,
			},
		},
	}
)

func TestMenuRenderer(t *testing.T) {
	xmlFilePath := filepath.Join("testdata", "menu.xml")
	xmlData, err := ioutil.ReadFile(xmlFilePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
		return
	}

	xmlString := string(xmlData)
	xml, err := menu.RenderXML()
	if err != nil {
		t.Fatalf("Expected error not to have occured :%v", err)
	}

	if processString(xml) != processString(xmlString) {
		t.Fatal("Expected xml to be equal to testdata")
	}
}

// TODO: Find a better alternative to this
func processString(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "\n", "")

	return str
}

// TODO: Find a better function name
func TestRendererWithJson(t *testing.T) {
	str, err := menu.RenderJSON()
	if err != nil {
		t.Fatalf("Expected error not to occur :%v", err)
	}

	t.Log(str)

}
