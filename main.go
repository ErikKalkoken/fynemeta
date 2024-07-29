package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Current version need to be injected via ldflags
var Version = "?"

func main() {
	flag.Usage = myUsage
	pathFlag := flag.String(
		"p",
		"",
		"path to the value in the format <key1>.<key2>\n"+
			"<key> can be the name of a key in a key/value pair, the name of a table\n"+
			"or the index of an array element (starting at 0). This option is mandatory.\n",
	)
	versionFlag := flag.Bool("v", false, "show the current version")
	flag.Parse()
	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	if *pathFlag == "" {
		exitWithError("Must provide path")
	}
	if len(flag.Args()) == 0 {
		exitWithError("Must provide file")
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
	v, err := findKey(data, p)
	if err != nil {
		exitWithError(err.Error())
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
	s := "Usage: tomlq [options] <file>:\n\n" +
		"Prints a value from a TOML file to stdout.\n" +
		"For more information please also see: https://github.com/ErikKalkoken/tomlq\n\n" +
		"Options:\n"
	fmt.Fprint(flag.CommandLine.Output(), s)
	flag.PrintDefaults()
}

func exitWithError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}

func findKey(data any, p []string) (any, error) {
	for i, k := range p {
		var v any
		var ok bool
		current := strings.Join(p[:i], ".")
		if current == "" {
			current = "<root>"
		}
		switch x := data.(type) {
		case map[string]any:
			v, ok = x[k]
			if !ok {
				return nil, fmt.Errorf("key \"%s\" not valid at: %s", k, current)
			}
		case []any:
			i, err := strconv.Atoi(k)
			if err != nil {
				return nil, fmt.Errorf("\"%s\" must be an integer at: %s", k, current)
			}
			if i > len(x)-1 {
				return nil, fmt.Errorf("%d is an invalid index at: %s", i, current)
			}
			v = x[i]
		case []map[string]any:
			i, err := strconv.Atoi(k)
			if err != nil {
				return nil, fmt.Errorf("\"%s\" must be an integer at: %s", k, current)
			}
			if i > len(x)-1 {
				return nil, fmt.Errorf("%d is an invalid index at: %s", i, current)
			}
			v = x[i]
		default:
			return nil, fmt.Errorf("value not found at: %s", current)
		}
		if i < len(p)-1 {
			data = v
		} else {
			switch x := reflect.ValueOf(v); x.Kind() {
			case reflect.Map, reflect.Slice:
				return nil, fmt.Errorf("\"%s\" must reference a simple data type at: %s", k, current)
			default:
				return v, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
