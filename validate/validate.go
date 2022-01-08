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
	SourceRanges map[string]hcl.Range
}

func NewResourceValidationInfo(attrs hcl.Attributes) *ResourceValidationInfo {
	result := map[string]hcl.Range{}
	for _, attr := range attrs {
		result[attr.Name] = attr.Range
	}
	return &ResourceValidationInfo{result}
}

func (info *ResourceValidationInfo) LogFatal(resourceId string, fieldName string, errorMessage string) {
	sourceRange := info.SourceRanges[fieldName]

	printToStdErr("Validation error when parsing resource %s: %s\n  on %s\n", resourceId, errorMessage, sourceRange)
	printLinesInRange(sourceRange)

	exitAndPrintStackTrace()
}

func LogFatalWithDiags(diags hcl.Diagnostics, format string, a ...interface{}) {
	printToStdErr(format, a...)

	for _, diag := range diags {
		if diag.Detail == "Unsuitable value: value must be known" {
			// useless diagnostic that always shows up
			continue
		}
		printToStdErr(diag.Error())
		if diag.Subject != nil {
			printLinesInRange(*diag.Subject)
		}
	}

	exitAndPrintStackTrace()
}

func LogFatalWithSourceRange(sourceRange hcl.Range, format string, a ...interface{}) {
	printToStdErr(format, a...)
	printToStdErr("  on %s\n", sourceRange)
	printLinesInRange(sourceRange)
	exitAndPrintStackTrace()
}

func LogInternalError(format string, a ...interface{}) {
	printToStdErr(format, a...)
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
			printToStdErr("  %s", line.String())
		}
		printToStdErr("")
	}
}

func printToStdErr(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format+"\n", a...)
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
			matchingLines = append(matchingLines, Line{
				LineNumber: scanner.Range().Start.Line,
				Content:    string(scanner.Bytes()),
			})
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return matchingLines, nil
}
