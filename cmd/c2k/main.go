package main

import (
	"c2k/internal/app/curl"
	"c2k/internal/app/kotlin"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: c2k curl ...\n")
		os.Exit(1)
	}

	request, err := curl.ParseCommand(os.Args[1:])

	if err != nil {
		log.Fatal(err)
	}

	ktFile, err := kotlin.GenAst(&request)

	if err != nil {
		log.Fatal(err)
	}

	err = kotlin.Serialize(os.Stdout, &ktFile)

	if err != nil {
		log.Fatal(err)
	}
}
