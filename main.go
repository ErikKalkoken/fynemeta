package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

var ErrNotFound = errors.New("not found")

func main() {
	flag.Usage = myUsage
	pathFlag := flag.String("p", "", "path to the value in the format key1.key2")
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
	process(filename, *pathFlag)
}

func process(filename string, path string) {
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
	switch x := v.(type) {
	case time.Time:
		fmt.Print(x.Format(time.RFC3339))
	default:
		fmt.Print(x)
	}
}

// myUsage writes a custom usage message to configured output stream.
func myUsage() {
	s := "Usage: tomlq -p <key1>[.<key2>[...]] <inputfile>:\n\n" +
		"Print a value from a TOML file.\n" +
		"For more information please also see: https://github.com/ErikKalkoken/tomlq\n\n"
	fmt.Fprint(flag.CommandLine.Output(), s)
}

func exitWithError(message string) {
	fmt.Printf("ERROR: %s\n", message)
	os.Exit(1)
}

func findKey(data any, p []string) (any, bool) {
	for i, k := range p {
		var v any
		var ok bool
		switch x := data.(type) {
		case map[string]any:
			v, ok = x[k]
			if !ok {
				return nil, false
			}
		case []any:
			i, err := strconv.Atoi(k)
			if err != nil {
				return nil, false
			}
			v = x[i]
		case []map[string]any:
			i, err := strconv.Atoi(k)
			if err != nil {
				return nil, false
			}
			v = x[i]
		default:
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
