package main

import (
	"c2k/internal/app/curl"
	"c2k/internal/app/kotlin"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

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

	var parseErr *curl.ParseError
	if errors.As(err, &parseErr) {
		switch parseErr.Kind {
		case curl.WarningKind:
			if len(parseErr.UnexpectedOptions) > 0 {
				_, _ = fmt.Fprintf(os.Stderr, "Skipping unexpected option[s]: %s\n", strings.Join(parseErr.UnexpectedOptions, " "))
			}
		case curl.NoCurlCommandKind:
			_, _ = fmt.Fprintf(os.Stderr, "Expected curl command, got %s\n", strings.Join(os.Args[1:], " "))
			os.Exit(1)
		case curl.NoUrlKind:
			_, _ = fmt.Fprintf(os.Stderr, "Expected URL argument for the curl command\n")
			os.Exit(1)
		case curl.UnknownParseErrorKind:
			_, _ = fmt.Fprintf(os.Stderr, "Unknown error occured while parsing the curl command\n")
			os.Exit(1)
		}
	}

	ktFile := kotlin.GenAst(command)

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
