package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	p := t.TempDir()
	x := `
	[alpha]
	bravo=5
	charlie=1979-05-27T00:32:00-07:00
	`
	fp := filepath.Join(p, "temp.toml")
	if err := os.WriteFile(fp, []byte(x), 0644); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		path string
		want string
	}{
		{"alpha.bravo", "5"},
		{"alpha.charlie", "1979-05-27T00:32:00-07:00"},
	}
	for _, tc := range cases {
		got, err := captureOutput(func() {
			lookup(fp, tc.path)
		})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tc.want, got)
	}
}

func captureOutput(f func()) (string, error) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer r.Close()
	os.Stdout = w
	f()
	os.Stdout = orig
	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func TestFindKey(t *testing.T) {
	x := `
	yankee="green"

	foxtrot=[1, 2, 3]

	india = 2023-05-27

	[[golf]]
	hotel=11

	[[golf]]
	hotel=22

	[alpha]
	bravo=4

	[charlie.delta]
	echo=42
	`
	var data any
	_, err := toml.Decode(x, &data)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		path []string
		ok   bool
		want any
	}{
		{[]string{"foxtrot"}, false, nil},
		{[]string{"yankee"}, true, "green"},
		{[]string{"alpha", "bravo"}, true, int64(4)},
		{[]string{"charlie", "delta", "echo"}, true, int64(42)},
		{[]string{"foxtrot", "1"}, true, int64(2)},
		{[]string{"golf", "1", "hotel"}, true, int64(22)},
		{[]string{"invalid", "1"}, false, nil},
		{[]string{"alpha", "invalid"}, false, nil},
		{[]string{"foxtrot", "4"}, false, nil},
		{[]string{"golf", "3", "hotel"}, false, nil},
		{[]string{"alpha"}, false, nil},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("should find key for %s", tc.path), func(t *testing.T) {
			got, err := findKey(data, tc.path)
			if tc.ok {
				if assert.NoError(t, err) {
					assert.Equal(t, tc.want, got)
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}
