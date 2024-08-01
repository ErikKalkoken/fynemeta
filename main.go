package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	filename = "FyneApp.toml"
)

// Current version need to be injected via ldflags
var Version = "0.0.0"

func main() {
	lookupCmd := flag.NewFlagSet("lookup", flag.ExitOnError)
	keyFlag := lookupCmd.String(
		"k",
		"",
		"Key path to the value in the format <key1>.<key2>\n"+
			"<key> can be the name of a key in a key/value pair, the name of a table\n"+
			"or the index of an array element (starting at 0). This parameter is mandatory.\n",
	)
	pathFlag := lookupCmd.String(
		"p",
		".",
		"Directory path to FyneApp.toml.\n",
	)
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}
	switch cmd := os.Args[1]; cmd {
	case "lookup":
		lookupCmd.Parse(os.Args[2:])
		if *keyFlag == "" {
			exitWithError("Must provide path")
		}
		path := filepath.Join(*pathFlag, filename)
		process(path, *keyFlag)
	case "version":
		versionCmd.Parse(os.Args[2:])
		fmt.Println(Version)
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	s := "Usage: fynemeta <command> [arguments]:\n\n" +
		"Use Fyne metadata in the build process.\n" +
		"For more information please also see: https://github.com/ErikKalkoken/tomlq\n\n" +
		"The commands are:\n" +
		"\tlookup\t\tprint a value from a TOML file to stdout\n" +
		"\tversion\t\tprint the tool's version\n\n" +
		"Use fynemeta <command> -h for more information about a command.\n"
	fmt.Print(s)
}

func process(path string, keys string) {
	text, err := os.ReadFile(path)
	if err != nil {
		exitWithError(err.Error())
	}
	var data any
	if _, err := toml.Decode(string(text), &data); err != nil {
		exitWithError("failed to decode file as TOML")
	}
	p := strings.Split(keys, ".")
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

func exitWithError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}

func findKey(data any, keys []string) (any, error) {
	for i, k := range keys {
		var v any
		var ok bool
		current := strings.Join(keys[:i], ".")
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
		if i < len(keys)-1 {
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
