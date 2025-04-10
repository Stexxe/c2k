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

		lines := strings.Split(string(b), "\n")
		var curlCommand strings.Builder

		sep := ""
		var expCode string

		for i, l := range lines {
			runes := []rune(l)

			if len(runes) > 1 && runes[0] == '/' && runes[1] == '/' {
				if strings.TrimSpace(string(runes[2:])) == "Dependencies:" {
					expCode = strings.Join(lines[i:], "\n")
					break
				}

				curlCommand.WriteString(sep)

				startOffset := 2
				if i == 0 {
					startOffset = 3
				}

				curlCommand.WriteString(string(runes[startOffset:]))

				sep = "\n"
			} else {
				expCode = strings.Join(lines[i:], "\n")
				break
			}
		}

		if curlCommand.Len() == 0 {
			t.Fatalf("Expected comment with the curl command on the first line, got %s ...", string(b[:16]))
		}

		cmdParsed := parseCurlCommand(curlCommand.String())

		command, err := curl.ParseCommand(cmdParsed)

		if err != nil {
			t.Fatal(err)
		}

		ktFile, err := kotlin.GenAst(command)

		if err != nil {
			t.Fatalf("%s: unexpected error %s", e.Name(), err)
		}

		var actual strings.Builder
		err = kotlin.Serialize(&actual, ktFile)

		if err != nil {
			t.Fatal(err)
		}

		diff := utils.Diff(
			fmt.Sprintf("%s-expected", e.Name()),
			[]byte(expCode),
			fmt.Sprintf("%s-actual", e.Name()),
			[]byte(actual.String()),
		)

		if len(diff) != 0 {
			t.Fatalf("%s: %s\n\n---Actual---\n%s", e.Name(), diff, actual.String())
		}
	}
}

func parseCurlCommand(cmd string) (args []string) {
	var argBuilder strings.Builder
	inQuote := false
	var quoteSym rune

	runes := []rune(cmd)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
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
				nextRune := runes[i+1]

				if nextRune == '\'' || nextRune == '"' || nextRune == '\\' {
					argBuilder.WriteRune(nextRune)
					i += 1
				} else {
					argBuilder.WriteRune(r)
				}
			} else {
				argBuilder.WriteRune(r)
			}
		}
	}

	if argBuilder.Len() > 0 {
		args = append(args, argBuilder.String())
	}

	return
}
