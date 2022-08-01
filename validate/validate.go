package validate

import (
	"bufio"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"golang.org/x/exp/constraints"
	"io/ioutil"
	"regexp"
)

type Validator interface {
	Check(value string, valueType interface{}) error
}

type RegexpValidator struct {
	pattern       string
	errorTemplate string
	regex         *regexp.Regexp
}

// Check validates provided string with a regexp based on the pattern and returns optional error.
func (r *RegexpValidator) Check(value string, valueType interface{}) error {
	r.regex = regexp.MustCompile(r.pattern)
	if !r.regex.MatchString(value) {
		return fmt.Errorf(r.errorTemplate, valueType)
	}
	return nil
}

// wordWithDotHyphenUnder80Pattern is a regexp pattern that matches string that contain alphanumerics, underscores, periods,
// and hyphens that start with alphanumeric and End alphanumeric or underscore. Limits size to 1-80.
// Based on https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
const wordWithDotHyphenUnder80Pattern = string(`^[a-zA-Z\d]$|^[a-zA-Z\d][\w\-.]{0,78}\w$`)

//NewWordWithDotHyphenUnder80Validator creates new RegexpValidator validating with wordWithDotHyphenUnder80Pattern.
func NewWordWithDotHyphenUnder80Validator() Validator {
	return &RegexpValidator{wordWithDotHyphenUnder80Pattern, "%s can contain only alphanumerics, underscores, periods, and hyphens;" +
		" must start with alphanumeric and end with alphanumeric or underscore and have 1-80 length", nil}
}

// cidrIPv4Pattern defines CIDR IPv4 notation with or without mask.
const cidrIPv4Pattern = string(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)([\/][0-3][0-2]?|[\/][1-2][0-9]|[\/][0-9])?$`)

//NewCIDRIPv4Check creates new RegexpValidator validating CIDR IPv4
func NewCIDRIPv4Check() Validator {
	return &RegexpValidator{cidrIPv4Pattern, "%s not valid CIDR IPv4 value", nil}
}

// matchWholeWordsPattern creates OR words matching regexp pattern with words. Regexp special characters must be
// escaped.
func matchWholeWordsPattern(words []string) string {
	var pattern string
	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		pattern += fmt.Sprintf(`^(%s)$`, word)
		if i != len(words)-1 {
			pattern += `|`
		}
	}
	return pattern
}

// NewProtocolCheck checks if provided protocol value is allowed in every deployment environment.
func NewProtocolCheck() Validator {
	return &RegexpValidator{matchWholeWordsPattern([]string{"tcp", "udp", "icmp", "\\*"}),
		"%s didn't match any protocol allowed value", nil}
}

// InRangeIncludingCheck represents <lowerBound, upperBound> range.
type InRangeIncludingCheck[T constraints.Ordered] struct {
	errorTemplate string
	lowerBound    T
	upperBound    T
}

func (i *InRangeIncludingCheck[T]) Check(value T, valueType interface{}) error {
	if value < i.lowerBound {
		return fmt.Errorf(i.errorTemplate, valueType, value, "lower", i.lowerBound)
	} else if value > i.upperBound {
		return fmt.Errorf(i.errorTemplate, valueType, value, "higher", i.lowerBound)
	}
	return nil
}

func newInRangeExcludingCheck[T constraints.Ordered](errorTemplate string, lower, upper T) InRangeIncludingCheck[T] {
	return InRangeIncludingCheck[T]{errorTemplate, lower, upper}
}

// NewPortCheck creates InRangeIncludingCheck that can validate port correctness.
func NewPortCheck() InRangeIncludingCheck[int32] {
	return newInRangeExcludingCheck[int32]("%v port %v cannot be %v than %v", 0, 65535)
}

// NewPriorityCheck creates InRangeIncludingCheck that can validate priority value.
func NewPriorityCheck() InRangeIncludingCheck[int64] {
	return newInRangeExcludingCheck[int64]("%v priority value %v cannot be %v than %v", 100, 4096)
}

type maxLengthValidator struct {
	maxLength int
}

func (m maxLengthValidator) Check(value string, valueType interface{}) error {
	if len(value) < 1 {
		return fmt.Errorf("%v cannot be empty", valueType)
	} else if len(value) > m.maxLength {
		return fmt.Errorf("%v maximum length is %v", valueType, m.maxLength)
	}
	return nil
}

type dummyValidator struct {
}

func (d dummyValidator) Check(value string, valueType interface{}) error {
	return nil
}

func NewDbUsernameValidator(engine resourcespb.DatabaseEngine, version string) Validator {
	switch engine {
	case resourcespb.DatabaseEngine_MYSQL:
		switch version {
		case "5.6":
			return &maxLengthValidator{16}
		default:
			return &maxLengthValidator{32}
		}
	case resourcespb.DatabaseEngine_POSTGRES:
		return &maxLengthValidator{63}
	case resourcespb.DatabaseEngine_MARIADB:
		return &maxLengthValidator{80}
	default:
		return &dummyValidator{}
	}
}

func NewDbPasswordValidator(engine resourcespb.DatabaseEngine) Validator {
	switch engine {
	case resourcespb.DatabaseEngine_MYSQL:
		// https://stackoverflow.com/a/31634299
		return &maxLengthValidator{32}
	case resourcespb.DatabaseEngine_POSTGRES:
		// https://stackoverflow.com/a/19499303
		return &maxLengthValidator{99}
	case resourcespb.DatabaseEngine_MARIADB:
		// https://github.com/MariaDB/server/blob/7c58e97bf6f80a251046c5b3e7bce826fe058bd6/mysys/get_password.c#L65
		return &maxLengthValidator{79}
	default:
		return &dummyValidator{}
	}
}

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
