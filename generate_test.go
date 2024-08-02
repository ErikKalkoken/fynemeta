package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("can deal with empty struct", func(t *testing.T) {
		_, err := generateXML(FyneApp{})
		assert.Error(t, err)

	})
}
func TestGenerateErrorHandling(t *testing.T) {
	website := "website"
	releases := map[string]string{"Summary": "summary", "License": "license", "Description": "description"}
	details := AppDetails{ID: "id", Name: "name"}
	linuxAndBSD := &LinuxAndBSD{
		Categories: []string{"Utility"},
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
				Release:     map[string]string{"Summary": "summary", "Description": "description"},
				Details:     details,
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website:     website,
				Release:     map[string]string{"License": "license", "Description": "description"},
				Details:     details,
				LinuxAndBSD: linuxAndBSD,
			},
			errMissingRequiredParameter,
		},
		{
			FyneApp{
				Website:     website,
				Release:     map[string]string{"Summary": "summary", "License": "license"},
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
				Details:     AppDetails{Name: "id"},
				LinuxAndBSD: linuxAndBSD,
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
