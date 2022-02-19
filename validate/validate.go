package validate

import (
	"bufio"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"io/ioutil"
	"os"
	"runtime/debug"
)

type ResourceValidationInfo struct {
	SourceRanges  map[string]hcl.Range
	BlockDefRange hcl.Range
}

func NewResourceValidationInfoFromContent(content *hcl.BodyContent, definitionRange hcl.Range) *ResourceValidationInfo {
	result := map[string]hcl.Range{}
	for _, attr := range content.Attributes {
		result[attr.Name] = attr.Range
	}
	// TODO: also map blocks

	return &ResourceValidationInfo{result, definitionRange}
}

func (info *ResourceValidationInfo) LogFatal(resourceId string, fieldName string, errorMessage string) {
	if _, ok := info.SourceRanges[fieldName]; ok {
		sourceRange := info.SourceRanges[fieldName]
		printToStdErrLn(
			"Validation error when parsing resource %s: %s\n  on %s\n", resourceId, errorMessage, sourceRange,
		)
		printLinesInRange(sourceRange)
	} else {
		printToStdErrLn(
			"Validation error when parsing resource %s (%s): %s\n", resourceId, info.BlockDefRange, errorMessage,
		)
	}

	exitAndPrintStackTrace()
}

func LogFatalWithDiags(diags hcl.Diagnostics, format string, a ...any) {
	printToStdErrLn(format, a...)

	for _, diag := range diags {
		if diag.Detail == "Unsuitable value: value must be known" {
			// useless diagnostic that always shows up
			continue
		}
		printToStdErrLn(diag.Error())
		if diag.Subject != nil {
			printLinesInRange(*diag.Subject)
		}
	}
	os.Exit(1)
}

func LogFatalWithSourceRange(sourceRange hcl.Range, format string, a ...any) {
	printToStdErr("%s: ", sourceRange)
	printToStdErrLn(format, a...)
	printLinesInRange(sourceRange)
	os.Exit(1)
}

func LogInternalError(format string, a ...any) {
	printToStdErrLn(format, a...)
	exitAndPrintStackTrace()
}

func exitAndPrintStackTrace() {
	debug.PrintStack()
	os.Exit(1)
}

func printLinesInRange(sourceRange hcl.Range) {
	lines, err := ReadLinesForRange(sourceRange)
	if err == nil {
		for _, line := range lines {
			printToStdErrLn("  %s", line.String())
		}
		printToStdErrLn("")
	}
}

func printToStdErrLn(format string, a ...any) {
	_, err := fmt.Fprintf(os.Stderr, format+"\n", a...)
	if err != nil {
		panic(err)
	}
}

func printToStdErr(format string, a ...any) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		panic(err)
	}
}

type Line struct {
	LineNumber int
	Content    string
}

func (l Line) String() string {
	return fmt.Sprintf("%d: %s", l.LineNumber, l.Content)
}

func ReadLinesForRange(sourceRange hcl.Range) ([]Line, error) {
	bytes, err := ioutil.ReadFile(sourceRange.Filename)
	if err != nil {
		return nil, err
	}
	return ReadLines(sourceRange, bytes)
}

func ReadLines(sourceRange hcl.Range, bytes []byte) ([]Line, error) {
	var matchingLines []Line
	scanner := hcl.NewRangeScanner(bytes, sourceRange.Filename, bufio.ScanLines)
	for inProgress := true; inProgress; inProgress = scanner.Scan() {
		if scanner.Range().Overlaps(sourceRange) {
			matchingLines = append(
				matchingLines, Line{
					LineNumber: scanner.Range().Start.Line,
					Content:    string(scanner.Bytes()),
				},
			)
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return matchingLines, nil
}
