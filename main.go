package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

const (
	sourceFilename  = "FyneApp.toml"
	generateCmdName = "generate"
	lookupCmdName   = "lookup"
	versionCmdName  = "version"
)

// Current version need to be injected via ldflags
var Version = "0.0.0"

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}
	switch cmd := os.Args[1]; cmd {
	case generateCmdName:
		cmd := flag.NewFlagSet(generateCmdName, flag.ExitOnError)
		sourceFlag := cmd.String(
			"s",
			".",
			"Path to source file (FyneApp.toml).",
		)
		destFlag := cmd.String(
			"d",
			".",
			"Path to where the AppStream file will be created",
		)
		typeFlag := cmd.String(
			"t",
			"",
			"type of metadata file to create. Supported: \"AppStream\". MANDATORY",
		)
		cmd.Parse(os.Args[2:])
		if *typeFlag == "" {
			exitWithError("Must define which type to generate")
		}
		if !slices.Contains([]string{"AppStream"}, *typeFlag) {
			exitWithError("Invalid type: " + *typeFlag)
		}
		source := filepath.Join(*sourceFlag, sourceFilename)
		if err := generate(source, *destFlag, *typeFlag); err != nil {
			exitWithError(err.Error())
		}
	case lookupCmdName:
		cmd := flag.NewFlagSet(lookupCmdName, flag.ExitOnError)
		keyFlag := cmd.String(
			"k",
			"",
			"Key path to the value in the format <key1>.<key2>\n"+
				"<key> can be the name of a key in a key/value pair, the name of a table\n"+
				"or the index of an array element (starting at 0). MANDATORY.\n",
		)
		sourceFlag := cmd.String(
			"s",
			".",
			"Source path to (FyneApp.toml).",
		)
		cmd.Parse(os.Args[2:])
		if *keyFlag == "" {
			exitWithError("Must provide path")
		}
		path := filepath.Join(*sourceFlag, sourceFilename)
		if err := lookup(path, *keyFlag); err != nil {
			exitWithError(err.Error())
		}
	case versionCmdName:
		cmd := flag.NewFlagSet(versionCmdName, flag.ExitOnError)
		cmd.Parse(os.Args[2:])
		fmt.Println(Version)
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print("Usage: fynemeta <command> [arguments]:\n\n" +
		"A tool to help use Fyne metadata in the build process\n" +
		"For more information please also see: https://github.com/ErikKalkoken/fynemeta\n\n" +
		"The commands are:\n")

	m := []struct {
		command     string
		description string
	}{
		{generateCmdName, "generate an appstream metadata file from the Fyne metadata file"},
		{lookupCmdName, "print a value from a TOML file to stdout"},
		{versionCmdName, "print the tool's version"},
	}
	for _, r := range m {
		fmt.Printf("\t%-15s %s\n", r.command, r.description)
	}
	fmt.Print("\nUse fynemeta <command> -h for more information about a command.\n")
}

func exitWithError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}
