package cellcube

import (
	"encoding/json"
	"encoding/xml"
)

//TODO: Implement line breaks

const Header = `<?xml version="1.0" encoding="ISO-8859-1"?>` + "\n"
const DocType = `<!DOCTYPE pages SYSTEM "cellflash-1.3.dtd">` + "\n"

type Menu struct {
	XMLName xml.Name `json:"-" xml:"pages"`
	Desc    string   `json:"descr" xml:"descr,attr,omitempty"`
	Pages   []Page   `json:"pages" xml:"page"`
}

type Page struct {
	Tag       string     `json:"tag" xml:"tag,attr,omitempty"`
	Body      string     `json:"body" xml:",innerxml"`
	PageLinks []PageLink `json:"pagelinks" xml:"a"`
}

type PageLink struct {
	HREF string `json:"href" xml:"href,attr"`
	A    string `json:"achor" xml:",innerxml"`
}

func (m *Menu) RenderXML() (string, error) {
	out, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}

	blob := []byte(Header + DocType + string(out))

	return string(blob), nil
}

func (m *Menu) RenderJSON() (string, error) {
	out, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
