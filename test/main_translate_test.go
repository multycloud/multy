package test

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/multycloud/multy/decoder"
	"github.com/multycloud/multy/encoder"
	"github.com/multycloud/multy/parser"
	"github.com/multycloud/multy/validate"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type TestFiles struct {
	InputFile  []string
	OutputFile string
	Dir        string
}

func TestTranslate(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	allTests := map[string]*TestFiles{}

	root := "_configs"
	err := filepath.WalkDir(
		root, func(path string, info os.DirEntry, err error) error {
			if info.IsDir() || (filepath.Ext(path) != ".tf" && filepath.Ext(path) != ".hcl") || strings.HasPrefix(
				filepath.Base(path), ".",
			) {
				return nil
			}

			ext := filepath.Ext(path)
			base := filepath.Dir(path)
			if _, ok := allTests[base]; !ok {
				allTests[base] = &TestFiles{}
			}

			if ext == ".tf" {
				allTests[base].OutputFile = path
			} else {
				allTests[base].InputFile = append(allTests[base].InputFile, path)
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
				test(*testFile, t)
			},
		)
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
	if len(testFiles.InputFile) == 0 {
		t.Fatalf("No tf file found for input file %s", testFiles.InputFile)
	}
	if testFiles.OutputFile == "" {
		t.Fatalf("No hcl file found for expected file %s", testFiles.OutputFile)
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
	f, diags := hclP.ParseHCL([]byte(hclOutput), "generated_file")
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

	actualContentBlockPrinter := NewBlockPrinter(actualContent, []byte(hclOutput))
	actualBlocks, errRange, err := groupByTypeAndId(actualContent)
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
	expectedContentBlockPrinter := NewBlockPrinter(expectedContent, expectedOutput)
	expectedBlocks, errRange, err := groupByTypeAndId(expectedContent)
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
				t.Errorf("unexpected block:\n%s", actualContentBlockPrinter.PrintBlock(block))
				continue
			}
			compare(t, expectedBlocks[typ][id], expectedContentBlockPrinter, actualContentBlockPrinter, block)
		}
	}

	for typ, blocks := range expectedBlocks {
		for id, block := range blocks {
			if _, ok := actualBlocks[typ][id]; !ok {
				t.Errorf("missing block %s\nexpected (%s):\n%s", id, block.DefRange, expectedContentBlockPrinter.PrintBlock(block))
			}
		}
	}

}

func groupByTypeAndId(content *hcl.BodyContent) (map[string]map[string]*hcl.Block, *hcl.Range, error) {
	result := map[string]map[string]*hcl.Block{}
	for t, blocks := range content.Blocks.ByType() {
		if result[t] == nil {
			result[t] = map[string]*hcl.Block{}
		}
		for _, block := range blocks {
			uniqueName := strings.Join(block.Labels, ".")
			if _, ok := result[t][uniqueName]; ok {
				if t != "provider" {
					return nil, &block.DefRange, fmt.Errorf("duplicate resource %s", uniqueName)
				}
				// TODO: IMPORTANT - handle duplicate resources
			}
			result[t][uniqueName] = block
		}
	}

	return result, nil, nil
}

func compare(t *testing.T, expected *hcl.Block, expectedPrinter *BlockPrinter, actualPrinter *BlockPrinter, actual *hcl.Block) {
	failed := false
	expectedAttrs, _ := expected.Body.JustAttributes()
	actualAttrs, _ := actual.Body.JustAttributes()

	blockRange := actual.DefRange
	for _, attr := range actualAttrs {
		blockRange = hcl.RangeOver(blockRange, attr.Range)
	}

	for name, attr := range actualAttrs {
		if _, ok := expectedAttrs[name]; !ok {
			errorMessage := fmt.Sprintf("in resouce '%s'\n", strings.Join(actual.Labels, "."))
			errorMessage += actualPrinter.PrintBlock(actual)

			errorMessage += fmt.Sprintf("\nunexpected attribute '%s' \n", name)
			actualLines, err := validate.ReadLines(attr.Range, actualPrinter.rawContent)
			if err != nil {
				panic(err)
			}
			for _, line := range actualLines {
				errorMessage += line.String() + "\n"
			}
			t.Errorf(errorMessage)
			failed = true
			continue
		}

		if !cmp.Equal(
			attr, expectedAttrs[name], cmp.Comparer(
				func(a, b cty.Value) bool {
					return a.Equals(b).True()
				},
			),
			cmpopts.IgnoreUnexported(hcl.TraverseRoot{}, hcl.TraverseAttr{}, hcl.TraverseIndex{}, hcl.TraverseSplat{}),
			cmpopts.IgnoreTypes(hcl.Range{}),
		) {
			actualLines, err := validate.ReadLines(attr.Range, actualPrinter.rawContent)
			if err != nil {
				panic(err)
			}
			expectedLines, err := validate.ReadLinesForRange(expectedAttrs[name].Range)
			if err != nil {
				panic(err)
			}
			errorMessage := fmt.Sprintf(
				"[%s] different attribute values for attr %s in resouce '%s'\n", attr.Range, name,
				strings.Join(actual.Labels, "."),
			)
			errorMessage += "expected:"
			for _, line := range expectedLines {
				errorMessage += "\n" + line.String()
			}
			errorMessage += "\n"
			errorMessage += "actual:"
			for _, line := range actualLines {
				errorMessage += "\n" + line.String()
			}
			failed = true
			t.Errorf(errorMessage)
		}
	}

	for name, attr := range expectedAttrs {
		if _, ok := actualAttrs[name]; !ok {
			errorMessage := fmt.Sprintf(
				"\n[%s] missing attribute '%s' in resouce '%s' \n", attr.Range, name, strings.Join(actual.Labels, "."),
			)
			expectedLines, err := validate.ReadLinesForRange(attr.Range)
			if err != nil {
				panic(err)
			}
			errorMessage += "expected:\n"
			for _, line := range expectedLines {
				errorMessage += line.String() + "\n"
			}
			failed = true
			t.Errorf(errorMessage)
		}

	}

	// If all attributes are correct so far, we still need to check nested blocks. HCL doesn't allow us to do that
	// without a schema, so we'll just have to compare everything.
	if !failed && !cmp.Equal(
		actual, expected, cmp.Comparer(
			func(a, b cty.Value) bool {
				return a.Equals(b).True()
			},
		), cmpopts.IgnoreUnexported(
			hclsyntax.Body{}, hcl.TraverseRoot{}, hcl.TraverseAttr{}, hcl.TraverseIndex{}, hcl.TraverseSplat{},
		), cmpopts.IgnoreTypes(hcl.Range{}),
	) {
		t.Errorf("some nested blocks differ within this block,\nexpected:%sactual:%s\n",
			expectedPrinter.PrintBlock(expected), actualPrinter.PrintBlock(actual))
	}
}

type BlockPrinter struct {
	allRanges  map[*hcl.Block]hcl.Range
	rawContent []byte
}

func (b *BlockPrinter) PrintBlock(block *hcl.Block) string {
	fmt.Println(b.allRanges[block])
	actualLines, err := validate.ReadLines(b.allRanges[block], b.rawContent)
	if err != nil {
		panic(err)
	}
	message := ""
	for _, line := range actualLines {
		message += line.String() + "\n"
	}
	return message
}

func NewBlockPrinter(content *hcl.BodyContent, rawContent []byte) *BlockPrinter {
	allPositions := []hcl.Pos{
		{
			Line:   0,
			Column: 0,
			Byte:   len(rawContent),
		},
	}
	for _, attr := range content.Attributes {
		allPositions = append(allPositions, attr.Range.Start)
	}
	for _, block := range content.Blocks {
		allPositions = append(allPositions, block.DefRange.Start)
	}
	sort.Slice(allPositions, func(i, j int) bool {
		return allPositions[i].Byte < allPositions[j].Byte
	})
	result := map[*hcl.Block]hcl.Range{}
	for _, block := range content.Blocks {
		// TODO: replace with binary search
		i := 0
		for i = range allPositions {
			if allPositions[i] == block.DefRange.Start {
				break
			}
		}
		result[block] = hcl.Range{
			Filename: block.DefRange.Filename,
			Start:    block.DefRange.Start,
			End: hcl.Pos{
				Line:   allPositions[i+1].Line - 1,
				Column: 0,
				Byte:   allPositions[i+1].Byte - 1,
			},
		}
	}

	return &BlockPrinter{
		allRanges:  result,
		rawContent: rawContent,
	}

}
