package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateXML(t *testing.T) {
	t.Run("can deal with empty struct", func(t *testing.T) {
		_, err := generateXML(FyneApp{})
		assert.Error(t, err)

	})
}
func TestGenerateXMLErrorHandling(t *testing.T) {
	website := "website"
	releases := map[string]string{"License": "license", "Description": "description"}
	details := AppDetails{ID: "id", Name: "name"}
	linuxAndBSD := &LinuxAndBSD{
		Categories: []string{"Utility"},
		Comment:    "Comment",
	}

	cases := []struct {
		app FyneApp
		err error
	}{
		{
			FyneApp{
				Website:     "",
				Release:     releases,
				Details:     details,
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website:     website,
				Release:     map[string]string{"Description": "description"},
				Details:     details,
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website:     website,
				Release:     map[string]string{"License": "license"},
				Details:     details,
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website: website,
				Release: releases,
				Details: details,
			},
			errMissingRequiredTable,
		},
		{
			FyneApp{
				Website:     website,
				Release:     releases,
				Details:     AppDetails{Name: "name"},
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website:     website,
				Release:     releases,
				Details:     AppDetails{ID: "id"},
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website: website,
				Release: releases,
				Details: details,
				LinuxAndBSD: &LinuxAndBSD{
					GenericName: "dummy",
				},
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website: website,
				Release: releases,
				Details: details,
				LinuxAndBSD: &LinuxAndBSD{
					Categories: []string{"Utility"},
				},
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website:     website,
				Release:     releases,
				Details:     details,
				LinuxAndBSD: linuxAndBSD,
			},
			nil,
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("should detect incomplete metadata. No#%d", i+1), func(t *testing.T) {
			_, err := generateXML(tc.app)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tempDir := t.TempDir()
	x := `
	Website = "https://github.com/ErikKalkoken/janice"
	[Details]
		Icon = "icon.png"
		Name = "Janice"
		ID = "io.github.erikkalkoken.janice"
		Version = "0.2.3"
		Build = 2

	[Release]
		BuildName = "janice"
		Description = "<p>A desktop app for viewing large JSON files.</p>"
		License = "MIT"
		Screenshots = "https://cdn.imgpile.com/f/0IrYBjJ_xl.png"
		ContentRating = "oars-1.1"

	[LinuxAndBSD]
		GenericName = "JSON viewer"
		Categories = ["Utility"]
		Comment = "A desktop app for viewing large JSON files"
		Keywords = ["json", "viewer"]
	`
	source := filepath.Join(tempDir, "FyneApp.toml")
	if err := os.WriteFile(source, []byte(x), 0644); err != nil {
		t.Fatal(err)
	}
	err := generate(source, tempDir, "AppStream")
	if assert.NoError(t, err) {
		target := filepath.Join(tempDir, "io.github.erikkalkoken.janice.metainfo.xml")
		byt, err := os.ReadFile(target)
		if assert.NoError(t, err) {
			xml := string(byt)
			assert.Contains(t, xml, "<component type=\"desktop-application\">")
		}
	}
}
