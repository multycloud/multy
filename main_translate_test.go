package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"multy-go/decoder"
	"multy-go/encoder"
	"multy-go/parser"
	"multy-go/validate"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

func TestTranslate(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	allTests := map[string]*TestFiles{}

	root := "./test"
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if info.IsDir() || (filepath.Ext(path) != ".tf" && filepath.Ext(path) != ".hcl") || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		if _, ok := allTests[base]; !ok {
			allTests[base] = &TestFiles{}
		}

		if ext == ".tf" {
			allTests[base].OutputFile = path
		} else {
			allTests[base].InputFile = path
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, testFile := range allTests {
		t.Run(filepath.Base(filepath.Dir(testFile.InputFile)), func(t *testing.T) {
			test(*testFile, t)
		})
	}
}

var tfConfigFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "terraform",
		},
		{
			Type:       "provider",
			LabelNames: []string{"name"},
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
		{
			Type: "locals",
		},
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
		{
			Type:       "module",
			LabelNames: []string{"name"},
		},
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "data",
			LabelNames: []string{"type", "name"},
		},
		{
			Type: "moved",
		},
	},
}

func test(testFiles TestFiles, t *testing.T) {

	t.Logf("testing file %s", testFiles.InputFile)
	if testFiles.InputFile == "" {
		t.Errorf("No tf file found for input file %s", testFiles.InputFile)
	}
	if testFiles.OutputFile == "" {
		t.Errorf("No hcl file found for expected file %s", testFiles.OutputFile)
	}
	p := parser.Parser{}
	parsedConfig := p.Parse(testFiles.InputFile)
	r := decoder.Decode(parsedConfig)
	hclOutput := encoder.Encode(r)

	var lines []string
	for i, line := range strings.Split(hclOutput, "\n") {
		lines = append(lines, fmt.Sprintf("%d:%s", i+1, line))
	}
	t.Logf("output:\n%s", strings.Join(lines, "\n"))

	hclP := hclparse.NewParser()
	f, diags := hclP.ParseHCL([]byte(hclOutput), testFiles.InputFile)
	if diags != nil {
		t.Fatal(diags)
	}
	actualContent, diags := f.Body.Content(tfConfigFileSchema)
	if diags != nil {
		t.Fatal(diags)
	}

	expectedOutput, err := ioutil.ReadFile(testFiles.OutputFile)
	if err != nil {
		t.Fatal(err)
	}

	f, diags = hclP.ParseHCLFile(testFiles.OutputFile)
	if diags != nil {
		t.Fatal(diags)
	}
	expectedContent, diags := f.Body.Content(tfConfigFileSchema)
	if diags != nil {
		t.Fatal(diags)
	}

	actualBlocks, errRange, err := groupAllLabels(actualContent)
	if err != nil {
		errorMessage := "found error in generated file\n" + err.Error()
		actualLines, err := validate.ReadLines(*errRange, []byte(hclOutput))
		if err != nil {
			panic(err)
		}
		for _, line := range actualLines {
			errorMessage += line.String() + "\n"
		}
		t.Fatal(errorMessage)
	}
	expectedBlocks, errRange, err := groupAllLabels(expectedContent)
	if err != nil {
		errorMessage := fmt.Sprintf("[%s] found error in expected file\n", errRange) + err.Error()
		actualLines, err := validate.ReadLines(*errRange, expectedOutput)
		if err != nil {
			panic(err)
		}
		for _, line := range actualLines {
			errorMessage += line.String() + "\n"
		}
		t.Fatal(errorMessage)
	}
	for typ, blocks := range actualBlocks {
		for id, block := range blocks {
			if _, ok := expectedBlocks[typ][id]; !ok {
				attrs, diags := block.Body.JustAttributes()
				if diags.HasErrors() {
					panic(diags.Error())
				}
				errorMessage := "unexpected block\n"
				errorMessage += printBlock(block, attrs, []byte(hclOutput))
				t.Error(errorMessage)
				continue
			}
			compare(expectedBlocks[typ][id], block, t, hclOutput)
		}
	}

	for typ, blocks := range expectedBlocks {
		for id, block := range blocks {
			if _, ok := actualBlocks[typ][id]; !ok {
				attrs, diags := block.Body.JustAttributes()
				if diags.HasErrors() {
					panic(diags.Error())
				}
				errorMessage := fmt.Sprintf("missing block %s\nexpected (%s):\n", id, block.DefRange)
				errorMessage += printBlock(block, attrs, expectedOutput)
				t.Error(errorMessage)
			}
		}
	}

}

func groupAllLabels(content *hcl.BodyContent) (map[string]map[string]*hcl.Block, *hcl.Range, error) {
	result := map[string]map[string]*hcl.Block{}
	for t, blocks := range content.Blocks.ByType() {
		if result[t] == nil {
			result[t] = map[string]*hcl.Block{}
		}
		for _, block := range blocks {
			uniqueName := strings.Join(block.Labels, ".")
			if _, ok := result[t][uniqueName]; ok {
				if t != "provider" {
					return nil, &block.DefRange, fmt.Errorf("duplicate resource %s\n", uniqueName)
				}
				// TODO: IMPORTANT - handle duplicate resources
			}
			result[t][uniqueName] = block
		}
	}

	return result, nil, nil
}

func printBlock(b *hcl.Block, attrs hcl.Attributes, bytes []byte) string {
	blockRange := b.DefRange
	for _, attr := range attrs {
		blockRange = hcl.RangeOver(blockRange, attr.Range)
	}
	actualLines, err := validate.ReadLines(blockRange, bytes)
	if err != nil {
		panic(err)
	}
	message := ""
	for _, line := range actualLines {
		message += line.String() + "\n"
	}
	return message
}

func compare(expected *hcl.Block, actual *hcl.Block, t *testing.T, actualFile string) {
	expectedAttrs, _ := expected.Body.JustAttributes()
	actualAttrs, _ := actual.Body.JustAttributes()

	blockRange := actual.DefRange
	for _, attr := range actualAttrs {
		blockRange = hcl.RangeOver(blockRange, attr.Range)
	}

	for name, attr := range actualAttrs {
		if _, ok := expectedAttrs[name]; !ok {
			errorMessage := fmt.Sprintf("in resouce '%s'\n", strings.Join(actual.Labels, "."))
			errorMessage += printBlock(actual, actualAttrs, []byte(actualFile))

			errorMessage += fmt.Sprintf("\nunexpected attribute '%s' \n", name)
			actualLines, err := validate.ReadLines(attr.Range, []byte(actualFile))
			if err != nil {
				panic(err)
			}
			for _, line := range actualLines {
				errorMessage += line.String() + "\n"
			}
			t.Errorf(errorMessage)
			continue
		}

		if !cmp.Equal(attr, expectedAttrs[name], cmp.Comparer(func(a, b cty.Value) bool {
			return a.Equals(b).True()
		}), cmpopts.IgnoreUnexported(hcl.TraverseRoot{}, hcl.TraverseAttr{}, hcl.TraverseIndex{}, hcl.TraverseSplat{}), cmpopts.IgnoreTypes(hcl.Range{})) {
			actualLines, err := validate.ReadLines(attr.Range, []byte(actualFile))
			if err != nil {
				panic(err)
			}
			expectedLines, err := validate.ReadLinesForRange(expectedAttrs[name].Range)
			if err != nil {
				panic(err)
			}
			errorMessage := fmt.Sprintf("[%s] different attribute values for attr %s in resouce '%s'\n", attr.Range, name, strings.Join(actual.Labels, "."))
			errorMessage += "expected: \n"
			for _, line := range expectedLines {
				errorMessage += line.String()
			}
			errorMessage += "\n"
			errorMessage += "actual: \n"
			for _, line := range actualLines {
				errorMessage += line.String()
			}
			t.Errorf(errorMessage)
		}
	}

	for name, attr := range expectedAttrs {
		if _, ok := actualAttrs[name]; !ok {
			errorMessage := fmt.Sprintf("\n[%s] missing attribute '%s' in resouce '%s' \n", attr.Range, name, strings.Join(actual.Labels, "."))
			expectedLines, err := validate.ReadLinesForRange(attr.Range)
			if err != nil {
				panic(err)
			}
			errorMessage += "expected:\n"
			for _, line := range expectedLines {
				errorMessage += line.String() + "\n"
			}
			t.Errorf(errorMessage)
		}

	}
}
