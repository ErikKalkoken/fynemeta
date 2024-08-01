package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Component is the top level element of the AppStream metadata XML.
type Component struct {
	XMLName         xml.Name     `xml:"component"`
	Id              string       `xml:"id"`
	Name            string       `xml:"name"`
	Summary         string       `xml:"summary"`
	MetadataLicense string       `xml:"metadata_license"`
	ProjectLicense  string       `xml:"project_license"`
	Description     any          `xml:"description"`
	LaunchAble      Parameter    `xml:"launchable"`
	URL             Parameter    `xml:"url"`
	ContentRating   Parameter    `xml:"content_rating,omitempty"`
	Screenshots     []Screenshot `xml:"screenshots>screenshot,omitempty"`
	Categories      []string     `xml:"categories>category"`
	Keywords        *Keywords    `xml:"keywords,omitempty"`
}

type Parameter struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type Keywords struct {
	Keyword []string `xml:"keyword,omitempty"`
}

type Screenshot struct {
	XMLName xml.Name `xml:"screenshot"`
	Type    string   `xml:"type,attr"`
	Image   string   `xml:"image"`
}

// appstream write an AppStream metadata file from a Fyne metadata file.
func appstream(path string) error {
	text, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var app FyneApp
	if _, err := toml.Decode(string(text), &app); err != nil {
		return fmt.Errorf("failed to decode file as TOML: %w", err)
	}
	stream := Component{
		Id:              app.Details.ID,
		Name:            app.Details.Name,
		Summary:         app.Release["Summary"],
		MetadataLicense: app.Release["License"],
		ProjectLicense:  app.Release["License"],
		LaunchAble: Parameter{
			Type:  "desktop-id",
			Value: app.Details.ID + ".desktop",
		},
		URL: Parameter{
			Type:  "homepage",
			Value: app.Website,
		},
		ContentRating: Parameter{
			Type: app.Release["ContentRating"],
		},
		Categories: app.LinuxAndBSD.Categories,
	}

	str := app.Release["Screenshots"]
	if str != "" {
		urls := strings.Split(str, ",")
		sst := make([]Screenshot, len(urls))
		for i, u := range urls {
			sst[i] = Screenshot{Type: "default", Image: u}
		}
		stream.Screenshots = sst
	}

	stream.Description = struct {
		Value string `xml:",innerxml"`
	}{app.Release["Description"]}

	if len(app.LinuxAndBSD.Keywords) > 0 {
		stream.Keywords = &Keywords{app.LinuxAndBSD.Keywords}
	}

	out, err := xml.MarshalIndent(stream, " ", "  ")
	if err != nil {
		return err
	}
	out2 := xml.Header + string(out)
	filename := app.Details.ID + ".appdata.xml"
	if err := os.WriteFile(filename, []byte(out2), 0664); err != nil {
		return err
	}
	fmt.Printf("Created appstream metadata file: %s\n", filename)
	return nil
}
