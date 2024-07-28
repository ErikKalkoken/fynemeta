package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var ErrNotFound = errors.New("not found")

func main() {
	flag.Usage = myUsage
	pathFlag := flag.String("k", "", "key path in the format key1.key2")
	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	if *pathFlag == "" {
		exitWithError("Must provide key path")
	}
	if len(flag.Args()) == 0 {
		exitWithError("Must provide file path")
	}
	filename := flag.Arg(0)
	v := process(filename, *pathFlag)
	fmt.Println(v)
}

func process(filename string, path string) any {
	text, err := os.ReadFile(filename)
	if err != nil {
		exitWithError(err.Error())
	}
	var data any
	if _, err := toml.Decode(string(text), &data); err != nil {
		exitWithError("failed to decode file as TOML")
	}
	p := strings.Split(path, ".")
	v, ok := findKey(data, p)
	if !ok {
		exitWithError(fmt.Sprintf("Failed to find key with path: %s", path))
	}
	return v
}

// myUsage writes a custom usage message to configured output stream.
func myUsage() {
	s := "Usage: tomlq -k <key1>[.<key2>[...]] <inputfile>:\n\n" +
		"Extracts a value from a TOML file.\n"
		// "For more information please see: https://github.com/ErikKalkoken/stellaris-tool\n\n"
	fmt.Fprint(flag.CommandLine.Output(), s)
}

func exitWithError(message string) {
	fmt.Printf("ERROR: %s\n", message)
	os.Exit(1)
}

func findKey(data any, p []string) (any, bool) {
	for i, k := range p {
		x, ok := data.(map[string]any)
		if !ok {
			return nil, false
		}
		v, ok := x[k]
		if !ok {
			return nil, false
		}
		if i < len(p)-1 {
			data = v
		} else {
			return v, true
		}
	}
	return nil, false
}
