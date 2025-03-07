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
	entries, err := os.ReadDir(testDir)

	if err != nil {
		t.Fatal(err)
	}

	for _, e := range entries {
		entryPath := filepath.Join(testDir, e.Name())

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

		cmd := string(b[start:i])
		expContent := b[i+1:]
		cmdParsed := strings.Split(cmd, " ")

		request, err := curl.ParseCommand(cmdParsed)

		if err != nil {
			t.Fatal(err)
		}

		ktFile, err := kotlin.GenAst(&request)

		if err != nil {
			t.Fatal(err)
		}

		var actual strings.Builder
		err = kotlin.Serialize(&actual, &ktFile)

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
			t.Fatalf("%s", diff)
		}
	}
}
