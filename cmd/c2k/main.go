package main

import (
	"c2k/internal/app/curl"
	"c2k/internal/app/kotlin"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

// TODO: Verbose mode is superset of include mode
// TODO: Error reporting
// TODO: README and Github fields

var Version string

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: c2k curl ...\n")
		os.Exit(1)
	}

	if os.Args[1] == "-V" || os.Args[1] == "--version" || os.Args[1] == "version" {
		fmt.Printf("Version: %s\n", getVersion())
		os.Exit(0)
	}

	command, err := curl.ParseCommand(os.Args[1:])

	if err != nil {
		log.Fatal(err)
	}

	ktFile, err := kotlin.GenAst(command)

	if err != nil {
		log.Fatal(err)
	}

	err = kotlin.Serialize(os.Stdout, ktFile)
	fmt.Println()

	if err != nil {
		log.Fatal(err)
	}
}

func getVersion() string {
	if Version != "" {
		return strings.Trim(Version, "\n\r")
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return fmt.Sprintf("dev-%s", setting.Value[:7])
			}
		}
	}

	return "dev"
}
