package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var plan = flag.Bool("plan", false, "plan")

func TestPlan(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	allTests := map[string]string{}

	root := "./test"
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".tf" || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		if _, ok := allTests[base]; !ok {
			allTests[base] = path
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, outfile := range allTests {
		if *plan {
			t.Run(filepath.Base(filepath.Dir(outfile)), func(t *testing.T) {
				testPlan(outfile, t)
			})
		}
	}
}

func testPlan(outfile string, t *testing.T) {
	t.Parallel()
	outputFilePath := fmt.Sprintf("%s_output.tf", strings.TrimSuffix(outfile, ".tf"))

	// if terraform providers change but .terraform already exists
	// run `find . -type d -name '.terraform' -exec rm -r {} +`
	// TODO add to make test clean
	if _, err := os.Stat(fmt.Sprintf("%s/.terraform", filepath.Dir(outputFilePath))); os.IsNotExist(err) {
		cmd := exec.Command("terraform", "init")
		cmd.Dir = filepath.Dir(outputFilePath)
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err, string(output))
	}

	cmd := exec.Command("terraform", "plan")
	cmd.Dir = filepath.Dir(outputFilePath)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(output))
}
