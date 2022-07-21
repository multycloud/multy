package validate_test

import (
	"github.com/multycloud/multy/validate"
	"regexp"
	"testing"
)

// TestWordWithDotHyphenUnder80Pattern checks whether validate.WordWithDotHyphenUnder80Pattern matches
// expected expressions
func TestWordWithDotHyphenUnder80Pattern(t *testing.T) {
	testRegexp, err := regexp.Compile(validate.WordWithDotHyphenUnder80Pattern)
	if err != nil {
		t.Fatalf("Could not compile regex: %s", validate.WordWithDotHyphenUnder80Pattern)
	}

	shouldMatch := []string{
		"a",
		"9",
		"aZ9",
		"ThisIs67dots..................................................................._",
	}
	shouldntMatch := []string{
		"",
		"_",
		"<someName",
		"ThisIs68dots...................................................................._",
		"Maybe?inThe.Middle_",
	}

	for _, name := range shouldMatch {
		if !testRegexp.MatchString(name) {
			t.Errorf("%s should match %s, but didn't", validate.WordWithDotHyphenUnder80Pattern, name)
		}
	}
	for _, name := range shouldntMatch {
		if testRegexp.MatchString(name) {
			t.Errorf("%s shouldn't match %s, but did", validate.WordWithDotHyphenUnder80Pattern, name)
		}
	}
}
