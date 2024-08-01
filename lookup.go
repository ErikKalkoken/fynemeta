package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func lookup(path string, keys string) error {
	text, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var data any
	if _, err := toml.Decode(string(text), &data); err != nil {
		return fmt.Errorf("failed to decode file as TOML: %w", err)
	}
	p := strings.Split(keys, ".")
	v, err := findKey(data, p)
	if err != nil {
		return err
	}
	switch x := v.(type) {
	case time.Time:
		fmt.Print(x.Format(time.RFC3339))
	default:
		fmt.Print(x)
	}
	return nil
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
