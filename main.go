package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	filename         = "FyneApp.toml"
	appstreamCmdName = "appstream"
	lookupCmdName    = "lookup"
	versionCmdName   = "version"
)

// Current version need to be injected via ldflags
var Version = "0.0.0"

func main() {
	appstreamCmd := flag.NewFlagSet(appstreamCmdName, flag.ExitOnError)
	pathFlag1 := appstreamCmd.String(
		"p",
		".",
		"Directory path to FyneApp.toml.\n",
	)

	lookupCmd := flag.NewFlagSet(lookupCmdName, flag.ExitOnError)
	keyFlag := lookupCmd.String(
		"k",
		"",
		"Key path to the value in the format <key1>.<key2>\n"+
			"<key> can be the name of a key in a key/value pair, the name of a table\n"+
			"or the index of an array element (starting at 0). This parameter is mandatory.\n",
	)
	pathFlag2 := lookupCmd.String(
		"p",
		".",
		"Directory path to FyneApp.toml.\n",
	)
	versionCmd := flag.NewFlagSet(versionCmdName, flag.ExitOnError)
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}
	switch cmd := os.Args[1]; cmd {
	case appstreamCmdName:
		appstreamCmd.Parse(os.Args[2:])
		path := filepath.Join(*pathFlag1, filename)
		if err := appstream(path); err != nil {
			exitWithError(err.Error())
		}
	case lookupCmdName:
		lookupCmd.Parse(os.Args[2:])
		if *keyFlag == "" {
			exitWithError("Must provide path")
		}
		path := filepath.Join(*pathFlag2, filename)
		if err := lookup(path, *keyFlag); err != nil {
			exitWithError(err.Error())
		}
	case versionCmdName:
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
	fmt.Print("Usage: fynemeta <command> [arguments]:\n\n" +
		"A tool to help use Fyne metadata in the build process\n" +
		"For more information please also see: https://github.com/ErikKalkoken/fynemeta\n\n" +
		"The commands are:\n")

	m := []struct {
		command     string
		description string
	}{
		{appstreamCmdName, "generate an appstream metadata file from the Fyne metadata file"},
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
