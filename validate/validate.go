package validate

import (
	"bufio"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"io/ioutil"
)

type ResourceValidationInfo struct {
	SourceRanges  map[string]hcl.Range
	BlockDefRange hcl.Range
	ResourceId    string
}

type ValidationError struct {
	ErrorMessage string
	ResourceId   string
	FieldName    string

	ResourceNotFound   bool
	ResourceNotFoundId string
}

func LogInternalError(format string, a ...any) {
	panic(fmt.Sprintf(format, a...))
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
