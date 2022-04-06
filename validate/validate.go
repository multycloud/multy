package validate

import (
	"bufio"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"io/ioutil"
	"os"
)

type ResourceValidationInfo struct {
	SourceRanges  map[string]hcl.Range
	BlockDefRange hcl.Range
	ResourceId    string
}

type ValidationError struct {
	SourceRange  hcl.Range
	ErrorMessage string
	ResourceId   string
	FieldName    string
}

func (e ValidationError) Print() {
	if !e.SourceRange.Empty() {
		printToStdErrLn("error when parsing resource %s: %s\n  on %s\n", e.ResourceId, e.ErrorMessage, e.SourceRange)
		printLinesInRange(e.SourceRange)
	} else {
		printToStdErrLn("error when parsing resource %s: %s\n", e.ResourceId, e.ErrorMessage)
	}
}

func (info *ResourceValidationInfo) NewError(errorMessage string, fieldName string) ValidationError {
	if _, ok := info.SourceRanges[fieldName]; ok {
		sourceRange := info.SourceRanges[fieldName]
		return ValidationError{
			SourceRange:  sourceRange,
			ErrorMessage: errorMessage,
			ResourceId:   info.ResourceId,
			FieldName:    fieldName,
		}
	}
	return ValidationError{
		ErrorMessage: errorMessage,
		ResourceId:   info.ResourceId,
		FieldName:    fieldName,
	}
}

func LogInternalError(format string, a ...any) {
	panic(fmt.Sprintf(format, a...))
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
