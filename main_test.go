package main

import (
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
	`
	fp := filepath.Join(p, "temp.toml")
	if err := os.WriteFile(fp, []byte(x), 0644); err != nil {
		t.Fatal(err)
	}
	got := process(fp, "alpha.bravo")
	assert.Equal(t, int64(5), got)
}

func TestFindKey(t *testing.T) {
	x := `
	yankee="green"
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
		want any
	}{
		{[]string{"yankee"}, "green"},
		{[]string{"alpha", "bravo"}, int64(4)},
		{[]string{"charlie", "delta", "echo"}, int64(42)},
	}
	for _, tc := range cases {
		got, ok := findKey(data, tc.path)
		if assert.True(t, ok) {
			assert.Equal(t, tc.want, got)
		}
	}

}
