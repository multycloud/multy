package test

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/proto/config"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestConfigFiles struct {
	InputFile  string
	OutputFile string
	Dir        string
}

func TestConfigTranslate(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	allTests := map[string]*TestConfigFiles{}

	root := "_configs"
	err := filepath.WalkDir(
		root, func(path string, info os.DirEntry, err error) error {
			if info.IsDir() || (!strings.HasSuffix(path, ".pb.json") && filepath.Ext(path) != ".tf") || strings.HasPrefix(
				filepath.Base(path), ".",
			) {
				return nil
			}

			ext := filepath.Ext(path)
			base := filepath.Dir(path)
			if _, ok := allTests[base]; !ok {
				allTests[base] = &TestConfigFiles{}
			}

			if ext == ".tf" {
				allTests[base].OutputFile = path
			} else {
				allTests[base].InputFile = path
			}
			return nil
		},
	)

	if err != nil {
		panic(err)
	}

	for dir, testFile := range allTests {
		t.Run(
			filepath.Base(dir), func(t *testing.T) {
				testConfig(*testFile, t)
			},
		)
	}
}

func testConfig(testFiles TestConfigFiles, t *testing.T) {
	if len(testFiles.InputFile) == 0 {
		// TODO: remove this return after migrating remaining tests
		return
		t.Fatalf("No textproto file found for input file %s", testFiles.OutputFile)
	}
	if testFiles.OutputFile == "" {
		t.Fatalf("No tf file found for expected file %s", testFiles.InputFile)
	}
	input, err := os.ReadFile(testFiles.InputFile)
	if err != nil {
		t.Fatalf("unable to open input file: %v", err)
	}

	c := config.Config{}
	err = jsonpb.UnmarshalString(string(input), &c)
	if err != nil {
		t.Fatalf("unable to parse input file: %v", err)
	}

	hclOutput, err := deploy.Translate(nil, &c, nil, nil)
	if err != nil && err != deploy.AwsCredsNotSetErr && err != deploy.AzureCredsNotSetErr {
		if s, ok := status.FromError(err); ok {
			fmt.Println(s.Details())
		}
		t.Fatalf("unable to translate: %v", err)
	}

	assertEqualHcl(t, []byte(hclOutput), testFiles.OutputFile)
}
