package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

var errMissingRequiredParameter = errors.New("missing required parameter")
var errMissingRequiredTable = errors.New("missing required table")

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

// generate write an AppStream metadata file from a Fyne metadata file.
func generate(source, destination, typ string) error {
	text, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	var app FyneApp
	if _, err := toml.Decode(string(text), &app); err != nil {
		return fmt.Errorf("failed to decode file as TOML: %w", err)
	}
	out, err := generateAppStreamData(app)
	if err != nil {
		return err
	}
	out2 := xml.Header + string(out)
	p := filepath.Join(destination, app.Details.ID+".appdata.xml")
	if err := os.WriteFile(p, []byte(out2), 0664); err != nil {
		return err
	}
	fmt.Printf("Created %s file: %s\n", typ, p)
	return nil
}

func generateAppStreamData(app FyneApp) ([]byte, error) {
	if app.Website == "" {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "Website")
	}
	if app.Release["Summary"] == "" {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "Release.Summary")
	}
	if app.Release["License"] == "" {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "Release.License")
	}
	if app.Release["Description"] == "" {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "Release.Description")
	}
	if app.Details.ID == "" {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "Details.ID")
	}
	if app.Details.Name == "" {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "Details.Name")
	}
	if app.LinuxAndBSD == nil {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredTable, "LinuxAndBSD")
	}
	if len(app.LinuxAndBSD.Categories) == 0 {
		return nil, fmt.Errorf("%w: %s", errMissingRequiredParameter, "LinuxAndBSD.Categories")
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

	return xml.MarshalIndent(stream, " ", "  ")
}
