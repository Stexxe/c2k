package e2e

import (
	"c2k/internal/app/curl"
	"c2k/internal/app/kotlin"
	"c2k/internal/app/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestConversion(t *testing.T) {
	testDir := "testData"
	casesDir := filepath.Join(testDir, "cases")
	entries, err := os.ReadDir(casesDir)

	if err != nil {
		t.Fatal(err)
	}

	for _, e := range entries {
		entryPath := filepath.Join(casesDir, e.Name())

		b, err := os.ReadFile(entryPath)

		if err != nil {
			t.Fatal(err)
		}

		if b[0] == '/' && b[1] == '/' {

		} else {
			t.Fatalf("Expected comment with the curl command on the first line, got %s ...", string(b[:16]))
		}

		i := 2
		for ; i < len(b) && b[i] == ' '; i++ {
		}
		start := i
		for ; i < len(b) && b[i] != '\n'; i++ {
		}

		expContent := b[i+1:]
		cmdParsed := parseCurlCommand(b[start:i])

		command, err := curl.ParseCommand(cmdParsed)

		if err != nil {
			t.Fatal(err)
		}

		ktFile, err := kotlin.GenAst(command)

		if err != nil {
			t.Fatal(err)
		}

		var actual strings.Builder
		err = kotlin.Serialize(&actual, ktFile)

		if err != nil {
			t.Fatal(err)
		}

		diff := utils.Diff(
			fmt.Sprintf("%s-expected", e.Name()),
			expContent,
			fmt.Sprintf("%s-actual", e.Name()),
			[]byte(actual.String()),
		)

		if len(diff) != 0 {
			t.Fatalf("%s: %s\n\n---Actual---\n%s", e.Name(), diff, actual.String())
		}
	}
}

func parseCurlCommand(cmd []byte) (args []string) {
	var argBuilder strings.Builder
	inQuote := false
	var quoteSym rune

	for i := 0; i < len(cmd); {
		r, sz := utf8.DecodeRune(cmd[i:])

		if inQuote {
			if r == quoteSym {
				inQuote = false
				args = append(args, argBuilder.String())
				argBuilder.Reset()
			} else {
				argBuilder.WriteRune(r)
			}
		} else {
			if r == ' ' {
				if argBuilder.Len() > 0 {
					args = append(args, argBuilder.String())
					argBuilder.Reset()
				}
			} else if r == '\'' || r == '"' {
				inQuote = true
				quoteSym = r
			} else if r == '\\' {
				nextRune, nextSize := utf8.DecodeRune(cmd[i+sz:])

				if nextRune == '\'' || nextRune == '"' || nextRune == '\\' {
					argBuilder.WriteRune(nextRune)
					i += nextSize
				} else {
					argBuilder.WriteRune(r)
				}
			} else {
				argBuilder.WriteRune(r)
			}
		}

		i += sz
	}

	if argBuilder.Len() > 0 {
		args = append(args, argBuilder.String())
	}

	return
}
